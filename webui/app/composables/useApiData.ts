export default function useApiData<T>(
    apiCall: () => Promise<{ data: T }>,
    key?: string
) {
    return useAsyncData(
        key || apiCall.name || Math.random().toString(36).substring(7),
        async () => {
            const response = await apiCall();
            return response.data;
        }
    );
};