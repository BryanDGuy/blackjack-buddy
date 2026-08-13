import { svelte } from '@sveltejs/vite-plugin-svelte';
import type { UserConfig } from 'vite';
import type { TestUserConfig } from 'vitest/config';

type Config = UserConfig & { test: TestUserConfig };

export default {
  plugins: [svelte()],
  build: {
    outDir: '../api/assets',
    emptyOutDir: true
  },
  resolve: {
    conditions: ['browser']
  },
  test: {
    environment: 'jsdom',
    globals: true
  }
} satisfies Config;
