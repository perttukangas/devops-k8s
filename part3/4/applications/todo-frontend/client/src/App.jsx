import { useState, useEffect } from 'react'

const App = () => {
  const [todos, setTodos] = useState([])
  const [newTodo, setNewTodo] = useState('')

  useEffect(() => {
    const fetchTodos = async () => {
      try {
        const response = await fetch('/todos')
        const data = await response.json()
        setTodos(data)
      } catch (error) {
        console.error('Error fetching todos:', error)
      }
    }

    fetchTodos()
  }, [])

  const handleCreateTodo = async () => {
    if (newTodo.trim() === '') return

    try {
      const response = await fetch('/todos', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ todo: newTodo }),
      })

      if (response.ok) {
        const createdTodo = await response.json()
        setTodos([...todos, createdTodo])
        setNewTodo('')
      } else {
        console.error('Error creating todo:', response.statusText)
      }
    } catch (error) {
      console.error('Error creating todo:', error)
    }
  }

  return (
    <>
      <h1>Hello from other branch</h1>
      <img
        src='/images/image.jpg'
        alt='Image description'
        style={{ maxWidth: '512px', maxHeight: '512px' }}
      />

      <h2>Todo List</h2>
      <input
        type='text'
        id='todoInput'
        maxLength='140'
        placeholder='Enter your todo'
        value={newTodo}
        onChange={(e) => setNewTodo(e.target.value)}
      />
      <button id='createTodo' onClick={handleCreateTodo}>
        Create
      </button>

      <ul id='todoList'>
        {todos.map((todo, index) => (
          <li key={index}>{todo.todo}</li>
        ))}
      </ul>
    </>
  )
}

export default App
