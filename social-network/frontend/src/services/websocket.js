class WebSocketService {
    constructor() {
        this.socket = null
        this.listeners = {}
        this.reconnectAttempts = 0
        this.maxReconnectAttempts = 5
    }

    connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const wsUrl = `${protocol}//${window.location.host}/ws`

        this.socket = new WebSocket(wsUrl)

        this.socket.onopen = () => {
            console.log('WebSocket connected')
            this.reconnectAttempts = 0
        }

        this.socket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data)
                if (this.listeners[data.type]) {
                    this.listeners[data.type](data.payload)
                }
            } catch (e) {
                console.error('Failed to parse WebSocket message', e)
            }
        }

        this.socket.onclose = () => {
            console.log('WebSocket disconnected')
            this.reconnect()
        }

        this.socket.onerror = (err) => {
            console.error('WebSocket error', err)
        }
    }

    reconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++
            const delay = 1000 * Math.pow(2, this.reconnectAttempts - 1)
            setTimeout(() => this.connect(), delay)
        }
    }

    send(type, to, payload) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify({ type, to, payload }))
        } else {
            console.warn('WebSocket not connected, message not sent')
        }
    }

    on(type, callback) {
        this.listeners[type] = callback
    }

    off(type) {
        delete this.listeners[type]
    }

    isConnected() {
        return this.socket && this.socket.readyState === WebSocket.OPEN
    }
}

export default new WebSocketService()