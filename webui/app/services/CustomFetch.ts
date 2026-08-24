export interface ApiBackendError {
  message: string;
}

type ExtractStatus<T> = T extends { status: infer S } ? S : number;

type ExtractData<T> = T extends { data: infer D } ? D : T;

export interface HttpResponse<T> {
  data: ExtractData<T>;
  status: ExtractStatus<T>;
  headers: Headers;
}

export const customFetch = async <T>(
  url: string,
  options: RequestInit
): Promise<any> => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase || ''
  const customOptions: RequestInit = {
    ...options,
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  };

  const response = await fetch(`${baseURL}${url}`, customOptions)

  if (!response.ok) {
    let errorMessage = `API Error: ${response.statusText}`;
    
    try {
      const errorPayload = (await response.json()) as ApiBackendError;
      if (errorPayload && errorPayload.message) {
        errorMessage = errorPayload.message;
      }
    } catch {}

    throw new Error(errorMessage);
  }

  if (response.status === 204) {
    return {
      data: {} as any,
      status: response.status as unknown as ExtractStatus<T>,
      headers: response.headers,
    };
  }

  const contentType = response.headers.get("content-type") || "";
  const isBlob = contentType.includes("image/") || contentType.includes("application/octet-stream");
  
  const rawData = isBlob ? await response.blob() : await response.json();

  return {
    data: rawData as unknown as ExtractData<T>,
    status: response.status as unknown as ExtractStatus<T>,
    headers: response.headers,
  };
}

export type ErrorType<Error> = Error;