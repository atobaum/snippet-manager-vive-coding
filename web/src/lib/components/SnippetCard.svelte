<script lang="ts">
	import { onMount, afterUpdate } from 'svelte';
	import hljs from 'highlight.js';

	export let snippet: any;
	export let onDelete: (name: string) => void;
	export let onEdit: (snippet: any) => void;
	export let onCopy: (text: string) => void;

	onMount(() => {
		highlightCode();
	});

	afterUpdate(() => {
		highlightCode();
	});

	function highlightCode() {
		setTimeout(() => {
			hljs.highlightAll();
		}, 100);
	}
</script>

<div class="bg-white rounded-xl shadow-lg border border-gray-100 hover:shadow-xl transition-all duration-200 transform hover:-translate-y-1">
	<div class="p-6">
		<div class="flex justify-between items-start mb-3">
			<div class="flex-1">
				<h3 class="font-bold text-gray-900 text-lg mb-1">{snippet.name}</h3>
				{#if snippet.language}
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
						{snippet.language}
					</span>
				{/if}
			</div>
			<button
				on:click={() => onDelete(snippet.name)}
				class="w-8 h-8 flex items-center justify-center text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg ml-3 flex-shrink-0 transition-colors duration-200"
				title="Delete snippet"
				aria-label="Delete snippet"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
				</svg>
			</button>
		</div>
		
		{#if snippet.description}
			<p class="text-gray-600 text-sm mb-4 leading-relaxed">{snippet.description}</p>
		{/if}
		
		{#if snippet.tags && snippet.tags.length > 0}
			<div class="flex flex-wrap gap-1 mb-4">
				{#each snippet.tags as tag}
					<span class="px-2 py-1 bg-emerald-100 text-emerald-800 text-xs rounded-md font-medium">{tag}</span>
				{/each}
			</div>
		{/if}
		
		<div class="bg-gray-900 rounded-lg p-4 mb-4 overflow-hidden">
			<pre class="text-sm text-gray-100 overflow-x-auto max-h-40 overflow-y-auto"><code class="language-{snippet.language || 'text'}">{snippet.command}</code></pre>
		</div>
		
		<div class="flex gap-2">
			<button
				on:click={() => onCopy(snippet.command)}
				class="flex-1 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors duration-200 flex items-center justify-center gap-2"
			>
				<svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
				</svg>
				Copy
			</button>
			<button
				on:click={() => onEdit(snippet)}
				class="w-10 h-10 bg-gray-100 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-200 transition-colors duration-200 flex items-center justify-center flex-shrink-0"
				title="Edit snippet"
				aria-label="Edit snippet"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
				</svg>
			</button>
		</div>
	</div>
</div>
