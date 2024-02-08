const APIClient = (method, data, resCallback) => $.ajax({
    url: '/places',
    type: method,
    dataType: 'json',
    success: (data) => {
        console.log(data);
    },
    error: (error) => {
        console.error('Error:', error);
    }
});

export default APIClient;