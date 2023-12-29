<template>
  <div v-if="isOpen" class="booking-modal">
    <div class="booking-content">
      <h2>Confirm Booking</h2>
      <p>Flight: {{ bookingDetails.routes.from }} - {{ bookingDetails.routes.destination }}</p>
      <p>Companies: {{ bookingDetails.companyNames.join(', ') }}</p>
      <p>Start Time: {{ this.formatDate(bookingDetails.startTime) }}</p>
      <p>Total price: {{ bookingDetails.totalPrice }}</p>
      <p>Total duration: {{ bookingDetails.totalDuration }}</p>

      <label for="firstName">First Name:</label>
      <input type="text" id="firstName" v-model="bookingDetails.firstName" />

      <label for="lastName">Last Name:</label>
      <input type="text" id="lastName" v-model="bookingDetails.lastName" />

      <button @click="this.sendDataAndConfirm">Confirm Booking</button>
    </div>
  </div>
</template>
<script>
export default {
  props: {
    isOpen: Boolean,
    bookingDetails: Object,
  },
  methods: {
    async sendDataAndConfirm() {
      const currentDateTime = new Date();
      const validUntilParts = this.bookingDetails.validUntil.split(/[/,:]+/);
      console.log(validUntilParts)
      const validUntilDateTime = new Date(
        validUntilParts[2],     
        validUntilParts[1] - 1,           
        validUntilParts[0],               
        validUntilParts[3],                
        validUntilParts[4],              
        validUntilParts[5]                
      );

      console.log(currentDateTime, validUntilDateTime);
      if (validUntilDateTime < currentDateTime) {
        alert("Unfortunately, this pricelist is outdated.");
        this.$router.go();
      } else if (!this.bookingDetails.firstName || !this.bookingDetails.lastName) {
        alert("Please enter both 'First Name' and 'Last Name' values.");
      } else {
        const requestData = { ...this.bookingDetails };
        delete requestData.validUntil;
        console.log("fetching");
        const response = await fetch(`http://localhost:8080/api/post`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(requestData),
        });
        if (response.status === 200) {
          this.$emit('confirm');
        } else if (response.status === 400) {
          this.$router.push({ name: "routeNotFound" });
        } else if (response.status === 500) {
          this.$router.push({ name: "internalError" });
        } else {
          console.log("Unhandled response status:", response.status);
        }
      }
    },
    formatDate(dateTimeString) {
      const parts = dateTimeString.split(/[-T:.Z]/);
      const date = new Date(
        Date.UTC(parts[0], parts[1] - 1, parts[2], parts[3], parts[4], parts[5])
      );
      return date.toLocaleString("en-GB");
    },
  },
};
</script>
  
<style scoped>
.booking-modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: white;
  padding: 20px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
}

.booking-content {
  text-align: center;
}

input {
  border: 1px solid black;
  padding: 2px;
}
</style>
  