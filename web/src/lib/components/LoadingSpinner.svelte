<script lang="ts">
	export let message: string = 'Loading...';
</script>

<div class="text-center py-8">
	<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
	<p class="mt-2 text-gray-600">{message}</p>
</div>
