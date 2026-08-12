import { defineConfig } from 'vitest/config';
import { svelte, vitePreprocess } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte({ preprocess: vitePreprocess() })],
  resolve: { conditions: ['browser'] },
  test: {
    environment: 'jsdom',
    globals: true
  }
});
