const APIClient = (method, data, resCallback) => $.ajax({
    url: 'http://localhost:3000/places',
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