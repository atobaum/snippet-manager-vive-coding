<script lang="ts">
	import { onMount } from 'svelte';
	import hljs from 'highlight.js';
	import 'highlight.js/styles/github.css';
	
	import Header from '$lib/components/Header.svelte';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import SnippetForm from '$lib/components/SnippetForm.svelte';
	import EditForm from '$lib/components/EditForm.svelte';
	import SnippetCard from '$lib/components/SnippetCard.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

	interface Snippet {
		name: string;
		description: string;
		language: string;
		tags: string[];
		command: string;
		created_at?: string;
		updated_at?: string;
	}

	let snippets: Snippet[] = [];
	let loading = true;
	let searchTerm = '';
	let showCreateForm = false;
	let editingSnippet: Snippet | null = null;
	let editSnippet = {
		name: '',
		description: '',
		language: '',
		command: '',
		tags: ''
	};
	let newSnippet = {
		name: '',
		description: '',
		language: '',
		command: '',
		tags: ''
	};

	onMount(async () => {
		await loadSnippets();
		// Initialize syntax highlighting after snippets are loaded
		setTimeout(() => {
			hljs.highlightAll();
		}, 100);
	});
	
	// Re-highlight when snippets change
	$: if (snippets.length > 0) {
		setTimeout(() => {
			hljs.highlightAll();
		}, 100);
	}

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
					language: newSnippet.language,
					command: newSnippet.command,
					tags: tags
				})
			});

			if (response.ok) {
				newSnippet = { name: '', description: '', language: '', command: '', tags: '' };
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

	function startEditSnippet(snippet: Snippet) {
		editingSnippet = snippet;
		editSnippet = {
			name: snippet.name,
			description: snippet.description,
			language: snippet.language,
			command: snippet.command,
			tags: snippet.tags.join(', ')
		};
		showCreateForm = false;
	}

	function cancelEdit() {
		editingSnippet = null;
		editSnippet = {
			name: '',
			description: '',
			language: '',
			command: '',
			tags: ''
		};
	}

	async function updateSnippet() {
		if (!editingSnippet) return;
		
		try {
			const tags = editSnippet.tags.split(',').map(tag => tag.trim()).filter(tag => tag !== '');
			const response = await fetch(`/api/snippets/${editingSnippet.name}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					description: editSnippet.description,
					language: editSnippet.language,
					command: editSnippet.command,
					tags: tags
				})
			});

			if (response.ok) {
				cancelEdit();
				await loadSnippets();
				// Re-highlight after update
				setTimeout(() => {
					hljs.highlightAll();
				}, 100);
			} else {
				alert('Failed to update snippet');
			}
		} catch (error) {
			console.error('Failed to update snippet:', error);
		}
	}

	function handleToggleCreateForm() {
		showCreateForm = !showCreateForm;
		if (showCreateForm) {
			cancelEdit();
		}
	}

	function handleCreateFormCancel() {
		showCreateForm = false;
	}

	$: filteredSnippets = snippets.filter(snippet =>
		snippet.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
		snippet.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
		snippet.tags.some(tag => tag.toLowerCase().includes(searchTerm.toLowerCase()))
	);
</script>

<div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-cyan-50">
	<div class="container mx-auto px-4 py-8 max-w-7xl">
		<Header />

		<SearchBar 
			bind:searchTerm={searchTerm}
			bind:showCreateForm={showCreateForm}
			onToggleCreateForm={handleToggleCreateForm}
		/>

		{#if showCreateForm}
			<SnippetForm 
				bind:newSnippet={newSnippet}
				onSubmit={createSnippet}
				onCancel={handleCreateFormCancel}
			/>
		{/if}

		{#if editingSnippet}
			<EditForm 
				bind:editingSnippet={editingSnippet}
				bind:editSnippet={editSnippet}
				onSubmit={updateSnippet}
				onCancel={cancelEdit}
			/>
		{/if}

		{#if loading}
			<LoadingSpinner message="Loading snippets..." />
		{:else if filteredSnippets.length === 0}
			<div class="text-center py-8">
				<p class="text-gray-500 text-lg">
					{searchTerm ? 'No snippets found matching your search.' : 'No snippets yet. Create your first one!'}
				</p>
			</div>
		{:else}
			<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
				{#each filteredSnippets as snippet}
					<SnippetCard 
						{snippet}
						onDelete={deleteSnippet}
						onEdit={startEditSnippet}
						onCopy={copyToClipboard}
					/>
				{/each}
			</div>
		{/if}
	</div>
</div>