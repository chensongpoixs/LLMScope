import { createApp } from 'vue'
import TestApp from './src/TestApp.vue'

const app = createApp(TestApp)
app.mount('#app')

console.log('Test app mounted')
