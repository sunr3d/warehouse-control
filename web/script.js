let currentToken = '';
let currentUser = '';
let currentRole = '';

// Проверить токен при загрузке
async function validateToken() {
    try {
        const response = await fetch('/items', {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });
        
        if (!response.ok) {
            // Токен невалиден, очистить
            logout();
            return false;
        }
        return true;
    } catch (error) {
        logout();
        return false;
    }
}

// Восстановить сессию при загрузке страницы
window.onload = async function() {
    const savedToken = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');
    const savedRole = localStorage.getItem('role');
    
    if (savedToken && savedUser && savedRole) {
        currentToken = savedToken;
        currentUser = savedUser;
        currentRole = savedRole;
        
        // Проверить валидность токена
        if (await validateToken()) {
            // Показать интерфейс
            document.getElementById('loginForm').style.display = 'none';
            document.getElementById('mainInterface').style.display = 'block';
            document.getElementById('currentUser').textContent = `${currentUser} (${currentRole})`;
            
            // Показать форму добавления для admin/manager
            if (currentRole === 'admin' || currentRole === 'manager') {
                document.getElementById('addItemForm').style.display = 'block';
            }
            
            loadItems();
        }
    }
};

// Логин
async function login() {
    const username = document.getElementById('userSelect').value;
    const password = document.getElementById('password').value;
    
    if (!username) {
        showError('Выберите пользователя');
        return;
    }

    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password })
        });

        if (response.ok) {
            const data = await response.json();
            currentToken = data.token;
            currentUser = data.username;
            currentRole = getRoleFromUsername(username);
            
            // Сохранить в localStorage
            localStorage.setItem('token', currentToken);
            localStorage.setItem('user', currentUser);
            localStorage.setItem('role', currentRole);
            
            document.getElementById('loginForm').style.display = 'none';
            document.getElementById('mainInterface').style.display = 'block';
            document.getElementById('currentUser').textContent = `${currentUser} (${currentRole})`;
            
            // Показываем форму добавления для admin и manager
            if (currentRole === 'admin' || currentRole === 'manager') {
                document.getElementById('addItemForm').style.display = 'block';
            }
            
            loadItems();
        } else {
            const error = await response.json();
            showError(error.error);
        }
    } catch (error) {
        showError('Ошибка подключения к серверу');
    }
}

// Выход
function logout() {
    currentToken = '';
    currentUser = '';
    currentRole = '';
    
    // Очистить localStorage
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('role');
    
    document.getElementById('loginForm').style.display = 'block';
    document.getElementById('mainInterface').style.display = 'none';
    document.getElementById('addItemForm').style.display = 'none';
}

// Получить роль по username
function getRoleFromUsername(username) {
    if (username === 'admin123') return 'admin';
    if (username === 'manager123') return 'manager';
    if (username === 'viewer123') return 'viewer';
    return 'unknown';
}

// Загрузить список товаров
async function loadItems() {
    try {
        const response = await fetch('/items', {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (response.ok) {
            const items = await response.json();
            displayItems(items);
        } else {
            showError('Ошибка загрузки товаров');
        }
    } catch (error) {
        showError('Ошибка подключения к серверу');
    }
}

// Отобразить товары
function displayItems(items) {
    const container = document.getElementById('itemsList');
    container.innerHTML = '';

    if (items.length === 0) {
        container.innerHTML = '<p>Товаров нет</p>';
        return;
    }

    const table = document.createElement('table');
    table.style.border = '1px solid black';
    table.style.borderCollapse = 'collapse';
    
    // Заголовки
    const header = table.insertRow();
    header.innerHTML = '<th>ID</th><th>Название</th><th>Описание</th><th>Количество</th><th>Создан</th><th>Обновлен</th><th>Действия</th>';
    
    // Строки товаров
    items.forEach(item => {
        const row = table.insertRow();
        row.innerHTML = `
            <td>${item.id}</td>
            <td>${item.name}</td>
            <td>${item.description || ''}</td>
            <td>${item.quantity}</td>
            <td>${new Date(item.created_at).toLocaleString('ru-RU')}</td>
            <td>${new Date(item.updated_at).toLocaleString('ru-RU')}</td>
            <td>
                ${currentRole === 'admin' || currentRole === 'manager' ? `
                    <button onclick="showHistory(${item.id}, '${item.name}')">История</button>
                ` : ''}
                ${currentRole === 'admin' || currentRole === 'manager' ? `
                    <button onclick="editItem(${item.id}, '${item.name}', '${item.description || ''}', ${item.quantity})">Изменить</button>
                ` : ''}
                ${currentRole === 'admin' ? `
                    <button onclick="deleteItem(${item.id})">Удалить</button>
                ` : ''}
            </td>
        `;
    });
    
    container.appendChild(table);
}

// Добавить товар
async function addItem() {
    const name = document.getElementById('itemName').value;
    const description = document.getElementById('itemDescription').value;
    const quantity = parseInt(document.getElementById('itemQuantity').value);

    if (!name || !quantity) {
        showError('Заполните название и количество');
        return;
    }

    try {
        const response = await fetch('/items', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentToken}`
            },
            body: JSON.stringify({ name, description, quantity })
        });

        if (response.ok) {
            // Очистить форму
            document.getElementById('itemName').value = '';
            document.getElementById('itemDescription').value = '';
            document.getElementById('itemQuantity').value = '';
            
            // Перезагрузить список
            loadItems();
        } else {
            const error = await response.json();
            showError(error.error);
        }
    } catch (error) {
        showError('Ошибка подключения к серверу');
    }
}

// Редактировать товар
async function editItem(id, currentName, currentDescription, currentQuantity) {
    const newName = prompt('Название:', currentName);
    if (newName === null) return;
    
    const newDescription = prompt('Описание:', currentDescription);
    if (newDescription === null) return;
    
    const newQuantity = prompt('Количество:', currentQuantity);
    if (newQuantity === null) return;
    
    const quantity = parseInt(newQuantity);
    if (isNaN(quantity) || quantity < 0) {
        showError('Некорректное количество');
        return;
    }

    try {
        const response = await fetch(`/items/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentToken}`
            },
            body: JSON.stringify({ name: newName, description: newDescription, quantity })
        });

        if (response.ok) {
            loadItems();
        } else {
            const error = await response.json();
            showError(error.error);
        }
    } catch (error) {
        showError('Ошибка подключения к серверу');
    }
}

// Удалить товар
async function deleteItem(id) {
    if (!confirm('Удалить товар?')) return;

    try {
        const response = await fetch(`/items/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (response.ok) {
            loadItems();
        } else {
            const error = await response.json();
            showError(error.error);
        }
    } catch (error) {
        showError('Ошибка подключения к серверу');
    }
}

// Показать историю
async function showHistory(itemId, itemName) {
    try {
        const response = await fetch(`/items/${itemId}/history`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });

        if (response.ok) {
            const data = await response.json();
            displayHistory(data, itemName);
        } else {
            const error = await response.json();
            showError(error.error);
        }
    } catch (error) {
        showError('Ошибка подключения к серверу');
    }
}

// Отобразить историю
function displayHistory(data, itemName) {
    document.getElementById('historyItemName').textContent = itemName;
    document.getElementById('historyModal').style.display = 'block';
    
    const container = document.getElementById('historyContent');
    container.innerHTML = '';

    if (data.items.length === 0) {
        container.innerHTML = '<p>История пуста</p>';
        return;
    }

    const table = document.createElement('table');
    table.style.border = '1px solid black';
    table.style.borderCollapse = 'collapse';
    
    // Заголовки
    const header = table.insertRow();
    header.innerHTML = '<th>Операция</th><th>Пользователь</th><th>Старое значение</th><th>Новое значение</th><th>Время</th>';
    
    // Строки истории
    data.items.forEach(entry => {
        const row = table.insertRow();
        row.innerHTML = `
            <td>${entry.operation}</td>
            <td>${entry.user_id}</td>
            <td>${entry.old_value || ''}</td>
            <td>${entry.new_value || ''}</td>
            <td>${new Date(entry.changed_at).toLocaleString('ru-RU')}</td>
        `;
    });
    
    container.appendChild(table);
}

// Закрыть историю
function closeHistory() {
    document.getElementById('historyModal').style.display = 'none';
}

// Показать ошибку
function showError(message) {
    const errorDiv = document.getElementById('loginError');
    errorDiv.textContent = message;
    setTimeout(() => {
        errorDiv.textContent = '';
    }, 5000);
}