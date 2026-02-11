import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig([
  {
    input: '../api/auth/v1/signup.yaml',
    output: 'lib/client/signup',
    plugins: [
      '@hey-api/client-next', 
      '@hey-api/typescript',
      '@hey-api/sdk',
    ],
  },
]);