<!-- CustomDropdown.vue -->
<template>
    <div class="custom-dropdown" ref="dropdown">
        <div class="selected-option" @click="toggleDropdown" :style="{ width: '180px' }">
            {{ selectedOption }}
        </div>
        <transition name="fade">
            <div v-if="isDropdownOpen" class="options-list">
                <div v-for="option in options" :key="option">
                    <label>
                        <input type="checkbox" :value="option" v-model="selectedOptions" @change="updateSelection" />
                        {{ option }}
                    </label>
                </div>
            </div>
        </transition>
    </div>
</template>
  
<script>
export default {
    props: {
        options: {
            type: Array,
            required: true,
        },
    },
    data() {
        return {
            isDropdownOpen: false,
            selectedOptions: [],
        };
    },
    computed: {
        selectedOption() {
            return this.selectedOptions.join(', ').slice(2);
        },
    },
    watch: {
        value: {
            immediate: true,
            handler(newValue) {
                this.selectedOptions = Array.isArray(newValue) ? [...newValue] : [newValue];
            },
        },
    },
    methods: {
        toggleDropdown() {
            this.isDropdownOpen = !this.isDropdownOpen;
            if (this.isDropdownOpen) {
                window.addEventListener('click', this.closeDropdownOnClickOutside);
            } else {
                window.removeEventListener('click', this.closeDropdownOnClickOutside);
            }
        },
        closeDropdownOnClickOutside(event) {
            if (!this.$refs.dropdown.contains(event.target)) {
                this.isDropdownOpen = false;
                window.removeEventListener('click', this.closeDropdownOnClickOutside);
            }
        },
        updateSelection() {
            this.$emit('input', this.selectedOptions);
        },
    },
};
</script>
  
<style scoped>
.custom-dropdown {
  position: relative;
  display: inline-block;
}

.selected-option {
  cursor: pointer;
  padding: 5px;
  border: 1px solid #ccc;
  height: 15px; /* Set the height as needed */
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.options-list {
  position: absolute;
  top: 100%;
  left: 0;
  z-index: 1;
  max-height: 150px;
  overflow-y: auto;
  border: 1px solid #ccc;
  background-color: #fff;
  width: 100%;
}

.options-list label {
  display: block;
  padding: 5px;
  cursor: pointer;
  height: 30px; /* Set the height as needed */
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}
</style>