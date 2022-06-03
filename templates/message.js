function btnClick() {
    const url = '/messages/send';
    const data = {message: document.message.message.value,
        idString: document.message.idString.value,
        platform: document.message.platform.value
    };

    try {
        const response = fetch(url, {
            method: 'POST', // или 'PUT'
            body: JSON.stringify(data), // данные могут быть 'строкой' или {объектом}!
            headers: {
                'Content-Type': 'application/json'
            }
        });
        console.log('Успех:');
    } catch (error) {
        console.error('Ошибка:', error);
    }
    location.reload();
}
