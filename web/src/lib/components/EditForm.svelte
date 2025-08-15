<script lang="ts">
	export let editingSnippet: any;
	export let editSnippet = {
		name: '',
		description: '',
		language: '',
		command: '',
		tags: ''
	};
	export let onSubmit: () => void;
	export let onCancel: () => void;
</script>

<div class="bg-white p-8 rounded-2xl shadow-xl mb-8 border border-gray-100 backdrop-blur-sm">
	<h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center gap-2">
		<svg class="w-6 h-6 text-green-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
		</svg>
		Edit Snippet: {editingSnippet?.name}
	</h2>
	<form on:submit|preventDefault={onSubmit} class="space-y-6">
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div>
				<label for="edit-snippet-description" class="block text-sm font-semibold text-gray-700 mb-2">Description</label>
				<input
					id="edit-snippet-description"
					type="text"
					bind:value={editSnippet.description}
					placeholder="What does this snippet do?"
					class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all duration-200 shadow-sm"
				/>
			</div>
			<div>
				<label for="edit-snippet-language" class="block text-sm font-semibold text-gray-700 mb-2">Language</label>
				<input
					id="edit-snippet-language"
					type="text"
					bind:value={editSnippet.language}
					placeholder="bash, python, javascript, go, etc."
					class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all duration-200 shadow-sm"
				/>
			</div>
		</div>
		<div>
			<label for="edit-snippet-tags" class="block text-sm font-semibold text-gray-700 mb-2">Tags</label>
			<input
				id="edit-snippet-tags"
				type="text"
				bind:value={editSnippet.tags}
				placeholder="utility, deployment, database, etc. (comma-separated)"
				class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all duration-200 shadow-sm"
			/>
		</div>
		<div>
			<label for="edit-snippet-command" class="block text-sm font-semibold text-gray-700 mb-2">Code/Command</label>
			<textarea
				id="edit-snippet-command"
				bind:value={editSnippet.command}
				rows="8"
				placeholder="Enter your code or command here..."
				class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all duration-200 shadow-sm font-mono text-sm resize-none"
			></textarea>
		</div>
		<div class="flex gap-3 pt-4">
			<button
				type="submit"
				class="flex-1 px-6 py-3 bg-gradient-to-r from-green-600 to-emerald-600 text-white font-medium rounded-xl hover:from-green-700 hover:to-emerald-700 transition-all duration-200 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 flex items-center justify-center gap-2"
			>
				<svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
				</svg>
				Update Snippet
			</button>
			<button
				type="button"
				on:click={onCancel}
				class="px-6 py-3 bg-gray-100 text-gray-700 font-medium rounded-xl hover:bg-gray-200 transition-all duration-200 flex items-center gap-2 whitespace-nowrap"
			>
				<svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
				Cancel
			</button>
		</div>
	</form>
</div>
