import js from '@eslint/js';
import tseslint from '@typescript-eslint/eslint-plugin';
import tsparser from '@typescript-eslint/parser';
import sveltePlugin from 'eslint-plugin-svelte';
import svelteParser from 'svelte-eslint-parser';

export default [
  js.configs.recommended,
  {
    files: ['**/*.ts', '**/*.js'],
    languageOptions: {
      parser: tsparser,
      parserOptions: {
        ecmaVersion: 2022,
        sourceType: 'module',
        project: './tsconfig.json'
      },
      globals: {
        console: 'readonly',
        window: 'readonly',
        document: 'readonly',
        fetch: 'readonly',
        HTMLElement: 'readonly',
        KeyboardEvent: 'readonly',
        Event: 'readonly'
      }
    },
    plugins: {
      '@typescript-eslint': tseslint
    },
    rules: {
      ...tseslint.configs.recommended.rules,
      '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
      '@typescript-eslint/no-explicit-any': 'error',
      'no-constant-condition': 'error',
      'no-empty': 'error',
      'no-implicit-coercion': 'error'
    }
  },
  {
    files: ['**/*.svelte'],
    languageOptions: {
      parser: svelteParser,
      parserOptions: {
        parser: tsparser,
        ecmaVersion: 2022,
        sourceType: 'module'
      },
      globals: {
        console: 'readonly',
        window: 'readonly',
        document: 'readonly',
        fetch: 'readonly',
        HTMLElement: 'readonly',
        KeyboardEvent: 'readonly',
        Event: 'readonly'
      }
    },
    plugins: {
      svelte: sveltePlugin,
      '@typescript-eslint': tseslint
    },
    rules: {
      ...sveltePlugin.configs.recommended.rules,
      ...tseslint.configs.recommended.rules,
      '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
      '@typescript-eslint/no-explicit-any': 'error',
      'no-constant-condition': 'error',
      'no-empty': 'error',
      'no-implicit-coercion': 'error',
      'svelte/no-at-html-tags': 'error',
      'svelte/valid-compile': 'error'
    }
  },
  {
    ignores: ['node_modules/**', 'dist/**', '*.config.js', '*.config.ts']
  }
];

