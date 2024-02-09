const APIClient = (endpoint, method, data, resCallback) => $.ajax({
    url: endpoint,
    type: method,
    data: data,
    dataType: 'json',
    success: (data) => {
        resCallback(data);
    },
    error: (error) => {
        console.error('Error:', error);
    }
});

export default APIClient;