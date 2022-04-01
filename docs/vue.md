# Vue

This doc outlines best practices for how we use Vue for Armeria. The purpose of this doc is to keep
the front-end code as consistent as possible. It is ordered such that each rule can be easily
referenced.

## 1. Components

### 1. General

- **a.** Components will always use the
  [Composition API](https://vuejs.org/guide/extras/composition-api-faq.html) over the
  [Options API](https://vuejs.org/guide/typescript/options-api.html).
- **b.** Components should favor `<script setup>` over `<script>` with a `setup()` function.
- **c.** Components should structure the composition as follows:
  - Imports
  - Props
  - Variable defaults
  - Computed store state/getters
  - Computed variables
  - Watches
  - Vue lifecycle hooks (e.g. `onMounted`)
  - Methods
- **d.** When defining props, always specify individual prop data types (e.g. `String`, `Number`,
  etc.).

### 2. Vuex Store Usage

- **a.** Place `const store = useStore()` directly above the computed constants which reference
  store state.

## 2. Javascript

- **a.** Always favor `async/await` over callback hell.
