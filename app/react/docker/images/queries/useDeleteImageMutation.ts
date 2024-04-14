import { RawAxiosRequestHeaders } from 'axios';
import { useMutation, useQueryClient } from 'react-query';

import { withInvalidate } from '@/react-tools/react-query';
import axios, { parseAxiosError } from '@/portainer/services/axios';
import { EnvironmentId } from '@/react/portainer/environments/types';

import { buildProxyUrl } from './build-url';
import { queryKeys } from './queryKeys';

export function useDeleteImageMutation(envId: EnvironmentId) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteImage,
    ...withInvalidate(queryClient, [queryKeys.base(envId)]),
  });
}

export async function deleteImage({
  environmentId,
  imageId,
  nodeName,
  force,
}: {
  environmentId: EnvironmentId;
  imageId: string;
  nodeName?: string;
  force?: boolean;
}) {
  const headers: RawAxiosRequestHeaders = {};

  if (nodeName) {
    headers['X-PortainerAgent-Target'] = nodeName;
  }

  try {
    await axios.delete(buildProxyUrl(environmentId, { id: imageId }), {
      headers,
      params: { force },
    });
  } catch (err) {
    throw parseAxiosError(err as Error, 'Unable to pull image');
  }
}
