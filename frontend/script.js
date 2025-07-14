document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('itemForm');
    const itemsContainer = document.getElementById('itemsList');

    // Сделаем функции глобальными
    window.updateItem = updateItem;
    window.deleteItem = deleteItem;

    // Загрузка элементов при старте
    loadItems();

    // Обработка добавления
    form.addEventListener('submit', e => {
        e.preventDefault();

        const newItem = {
            title: document.getElementById('title').value,
            description: document.getElementById('description').value
        };

        fetch('/items', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(newItem)
        })
        .then(response => response.json())
        .then(() => {
            form.reset();
            loadItems();
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to add item');
        });
    });

    // Загрузка элементов
    function loadItems() {
        itemsContainer.innerHTML = '<div class="loading">Loading items...</div>';
        
        fetch('/items')
            .then(response => {
                if (!response.ok) throw new Error('Network response was not ok');
                return response.json();
            })
            .then(items => {
                renderItems(items);
            })
            .catch(error => {
                console.error('Error:', error);
                itemsContainer.innerHTML = '<div class="error">Failed to load items</div>';
            });
    }
    
    function renderItems(items) {
        itemsContainer.innerHTML = '';
        
        if (items.length === 0) {
            itemsContainer.innerHTML = '<div class="loading">No items found</div>';
            return;
        }
        
        items.forEach(item => {
            const itemDiv = document.createElement('div');
            itemDiv.className = 'item';
            itemDiv.innerHTML = `
                <h3>${item.title}</h3>
                <p>${item.description || 'No description'}</p>
                <small>ID: ${item.id}</small>
                <div class="actions">
                    <button onclick="updateItem(${item.id})">Edit</button>
                    <button onclick="deleteItem(${item.id})">Delete</button>
                </div>
            `;
            itemsContainer.appendChild(itemDiv);
        });
    }
    
    function updateItem(id) {
        const newTitle = prompt("Enter new title:");
        if (newTitle === null) return; // Отмена
        
        const newDescription = prompt("Enter new description:") || "";

        fetch(`/update-item?id=${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                title: newTitle,
                description: newDescription
            })
        })
        .then(response => {
            if (!response.ok) throw new Error('Update failed');
            loadItems();
            alert('Item updated successfully!');
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to update item');
        });
    }

    function deleteItem(id) {
        if (!confirm("Are you sure you want to delete this item?")) return;
        
        fetch(`/delete-item?id=${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (!response.ok) throw new Error('Delete failed');
            loadItems();
            alert('Item deleted successfully!');
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to delete item');
        });
    }
});
