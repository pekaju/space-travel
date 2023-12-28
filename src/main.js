import { createApp } from 'vue';
import App from './App.vue';
import { createRouter, createWebHashHistory } from 'vue-router';
import LandingPage from './views/LandingPage.vue';
import ResultsPage from './views/ResultsPage.vue';
import RouteNotFound from './views/RouteNotFound.vue'
import InternalError from './views/InternalError.vue'

const routes = [
  {
    path: '/',
    name: 'landing',
    component: LandingPage,
  },
  {
    path: '/results',
    name: 'results',
    component: ResultsPage,
  },
  {
    path:'/routeNotFound',
    name: 'routeNotFound',
    component: RouteNotFound
  },
  {
    path: '/internalError',
    name: 'internalError',
    component: InternalError
  }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

const app = createApp(App);
app.use(router);

app.mount('#app');
