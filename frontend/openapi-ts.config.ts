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
  {
    input: '../api/media/v1/media.yaml',
    output: 'lib/client/media/v1/',
    plugins: [
      '@hey-api/client-next', 
      '@hey-api/typescript',
      '@hey-api/sdk',
    ],
  },
]);