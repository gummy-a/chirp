import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig([
  {
    input: '../api/auth/v1/auth.yaml',
    output: 'lib/client/auth/v1/',
    plugins: [
      '@hey-api/client-next', 
      '@hey-api/typescript',
      '@hey-api/sdk',
    ],
  },
]);