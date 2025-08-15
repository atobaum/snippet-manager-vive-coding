<script lang="ts">
	import { onMount } from 'svelte';
	
	interface Snippet {
		name: string;
		description: string;
		tags: string[];
		command: string;
		created_at?: string;
		updated_at?: string;
	}

	let snippets: Snippet[] = [];
	let loading = true;
	let searchTerm = '';
	let showCreateForm = false;
	let newSnippet = {
		name: '',
		description: '',
		command: '',
		tags: ''
	};

	onMount(async () => {
		await loadSnippets();
	});

	async function loadSnippets() {
		try {
			const response = await fetch('/api/snippets');
			snippets = await response.json();
		} catch (error) {
			console.error('Failed to load snippets:', error);
		} finally {
			loading = false;
		}
	}

	async function createSnippet() {
		try {
			const tags = newSnippet.tags.split(',').map(t => t.trim()).filter(t => t);
			const response = await fetch('/api/snippets', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					name: newSnippet.name,
					description: newSnippet.description,
					command: newSnippet.command,
					tags: tags
				})
			});

			if (response.ok) {
				newSnippet = { name: '', description: '', command: '', tags: '' };
				showCreateForm = false;
				await loadSnippets();
			} else {
				alert('Failed to create snippet');
			}
		} catch (error) {
			console.error('Failed to create snippet:', error);
		}
	}

	async function deleteSnippet(name: string) {
		if (!confirm(`Are you sure you want to delete "${name}"?`)) return;

		try {
			const response = await fetch(`/api/snippets/${name}`, {
				method: 'DELETE'
			});

			if (response.ok) {
				await loadSnippets();
			} else {
				alert('Failed to delete snippet');
			}
		} catch (error) {
			console.error('Failed to delete snippet:', error);
		}
	}

	function copyToClipboard(text: string) {
		navigator.clipboard.writeText(text);
		alert('Copied to clipboard!');
	}

	$: filteredSnippets = snippets.filter(snippet =>
		snippet.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
		snippet.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
		snippet.tags.some(tag => tag.toLowerCase().includes(searchTerm.toLowerCase()))
	);
</script>

<div class="container mx-auto px-4 py-8 max-w-6xl">
	<header class="mb-8">
		<h1 class="text-4xl font-bold text-gray-800 mb-4">📝 Snippet Manager</h1>
		<p class="text-gray-600">Manage your code snippets efficiently</p>
	</header>

	<!-- Search and Create Button -->
	<div class="flex flex-col sm:flex-row gap-4 mb-6">
		<div class="flex-1">
			<input
				type="text"
				bind:value={searchTerm}
				placeholder="Search snippets..."
				class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
			/>
		</div>
		<button
			on:click={() => showCreateForm = !showCreateForm}
			class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
		>
			{showCreateForm ? 'Cancel' : 'Create Snippet'}
		</button>
	</div>

	<!-- Create Form -->
	{#if showCreateForm}
		<div class="bg-white p-6 rounded-lg shadow-md mb-6 border">
			<h2 class="text-xl font-semibold mb-4">Create New Snippet</h2>
			<form on:submit|preventDefault={createSnippet} class="space-y-4">
				<div>
					<label for="snippet-name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
					<input
						id="snippet-name"
						type="text"
						bind:value={newSnippet.name}
						required
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					/>
				</div>
				<div>
					<label for="snippet-description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
					<input
						id="snippet-description"
						type="text"
						bind:value={newSnippet.description}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					/>
				</div>
				<div>
					<label for="snippet-tags" class="block text-sm font-medium text-gray-700 mb-1">Tags (comma-separated)</label>
					<input
						id="snippet-tags"
						type="text"
						bind:value={newSnippet.tags}
						placeholder="javascript, utility, function"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					/>
				</div>
				<div>
					<label for="snippet-command" class="block text-sm font-medium text-gray-700 mb-1">Command/Code</label>
					<textarea
						id="snippet-command"
						bind:value={newSnippet.command}
						rows="6"
						required
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
					></textarea>
				</div>
				<div class="flex gap-2">
					<button
						type="submit"
						class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
					>
						Create
					</button>
					<button
						type="button"
						on:click={() => showCreateForm = false}
						class="px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600 transition-colors"
					>
						Cancel
					</button>
				</div>
			</form>
		</div>
	{/if}

	<!-- Loading State -->
	{#if loading}
		<div class="text-center py-8">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<p class="mt-2 text-gray-600">Loading snippets...</p>
		</div>
	{:else if filteredSnippets.length === 0}
		<div class="text-center py-8">
			<p class="text-gray-500 text-lg">
				{searchTerm ? 'No snippets found matching your search.' : 'No snippets yet. Create your first one!'}
			</p>
		</div>
	{:else}
		<!-- Snippets Grid -->
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each filteredSnippets as snippet}
				<div class="bg-white rounded-lg shadow-md border hover:shadow-lg transition-shadow">
					<div class="p-4">
						<div class="flex justify-between items-start mb-2">
							<h3 class="font-semibold text-gray-800 truncate">{snippet.name}</h3>
							<button
								on:click={() => deleteSnippet(snippet.name)}
								class="text-red-500 hover:text-red-700 ml-2 flex-shrink-0"
								title="Delete snippet"
							>
								🗑️
							</button>
						</div>
						
						{#if snippet.description}
							<p class="text-gray-600 text-sm mb-3">{snippet.description}</p>
						{/if}
						
						{#if snippet.tags.length > 0}
							<div class="flex flex-wrap gap-1 mb-3">
								{#each snippet.tags as tag}
									<span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">{tag}</span>
								{/each}
							</div>
						{/if}
						
						<div class="bg-gray-50 rounded p-3 mb-3">
							<pre class="text-sm text-gray-800 whitespace-pre-wrap overflow-x-auto max-h-32 overflow-y-auto">{snippet.command}</pre>
						</div>
						
						<button
							on:click={() => copyToClipboard(snippet.command)}
							class="w-full px-3 py-2 bg-blue-600 text-white text-sm rounded hover:bg-blue-700 transition-colors"
						>
							📋 Copy to Clipboard
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
