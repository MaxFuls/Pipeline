let fieldCounter = 0;

function addField() {
    fieldCounter++;
    const container = document.getElementById('fieldsContainer');
    const fieldItem = document.createElement('div');
    fieldItem.className = 'field-item';
    fieldItem.innerHTML = `
        <input type="text" 
               placeholder="Введите дополнительную информацию..." 
               id="field_${fieldCounter}">
        <button class="btn btn-remove" onclick="removeField(this)">
            ✕
        </button>
    `;
    container.appendChild(fieldItem);
    updateFieldCount();
}

function removeField(button) {
    button.parentElement.remove();
    updateFieldCount();
}

function updateFieldCount() {
    const count = document.querySelectorAll('.field-item').length;
    document.getElementById('fieldCount').textContent = count;
}

async function submitForm() {
    const mainMessage = document.getElementById('mainMessage').value;
    const additionalFields = [];
    
    document.querySelectorAll('.field-item input').forEach(input => {
        if (input.value.trim()) {
            additionalFields.push(input.value);
        }
    });

    // Показываем индикатор загрузки
    showResponse('Отправка данных...', 'loading');

    try {
        const response = await fetch('/submit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                mainMessage: mainMessage,
                additionalFields: additionalFields
            })
        });

        const result = await response.json();
        
        if (result.success) {
            showResponse(`✅ ${result.message}\n\n📝 Основное сообщение: ${result.data.mainMessage}\n📋 Дополнительные поля: ${result.data.additionalFields.length} шт.\n⏰ Время: ${result.data.timestamp}`, 'success');
        } else {
            showResponse(`❌ Ошибка: ${result.message}`, 'error');
        }
    } catch (error) {
        showResponse(`❌ Ошибка сети: ${error.message}`, 'error');
    }
}

function showResponse(message, type) {
    const responseArea = document.getElementById('response');
    responseArea.textContent = message;
    responseArea.className = `response-area response-${type}`;
    responseArea.style.display = 'block';
    
    // Автоматически скрываем через 5 секунд для успешных ответов
    if (type === 'success') {
        setTimeout(() => {
            responseArea.style.display = 'none';
        }, 5000);
    }
}

// Добавляем первое поле при загрузке страницы
window.onload = function() {
    addField();
};