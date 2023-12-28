<template>
    <div class="search-form">
        <form @submit.prevent="searchFlights">
            <div class="input-outer-box">

                <div class="origin-outer-box">
                    <div class="origin-inner-box">
                        <label for="from">From</label>
                        <input v-model="from" type="text" id="from" placeholder="Choose origin planet..."/>
                    </div>
                </div>

                <div class="destination-outer-box">
                    <div class="destination-inner-box">
                        <label for="destination">Destination</label>
                        <input v-model="destination" type="text" id="destination" placeholder="Choose destination planet"/>
                    </div>
                </div>
                <button type="submit">Search Trips</button>
            </div>
        </form>
    </div>
</template>
<script>
export default {
    data() {
        return {
            from: "",
            destination: "",
        };
    },
    methods: {
        searchFlights() {
            if (this.from && this.destination) {
                this.$router.push({
                    name: 'results',
                    query: { from: this.from, destination: this.destination },
                });
            } else {
                alert("Please enter both 'From' and 'Destination' values.");
            }
        },
    },
}
</script>
<style>
.input-outer-box {
    display: flex;
    justify-content: space-around;
}

.origin-outer-box,
.destination-outer-box {
    height: 72px;
    width: 176px;
    background-color: white;
    border-radius: 8px;
    display: flex;
    align-items: center;
    margin: 10px;
}

.origin-outer-box:focus-within,
.destination-outer-box:focus-within {
    outline: 3px solid rgb(43, 93, 255);
}

input {
    border: none;
    padding-left: 10px;
    padding-top: 5px;
    text-overflow: ellipsis;
}

input:focus {
    outline: none;
}

input::placeholder {
    font-size: 17px;
}

label {
    padding-left: 10px;
    padding-bottom: 5px;
}

button {
    margin: 10px;
    border-radius: 8px;
    border: none;
    font-size: 17px;
}

button:hover {
    cursor: pointer;
}

datalist {
    max-height: 100px;
    /* Adjust as needed */
    overflow-y: auto;
    border: 1px solid #ccc;
    background-color: #fff;
}
</style>