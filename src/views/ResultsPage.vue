<template>
    <div class="search-form-area">
        <SearchForm></SearchForm>
    </div>
    <div class="big-div">
        <div class="info">
            <p>From: {{ from }}</p>
            <p>Destination: {{ destination }}</p>
            <p>Distance: {{ distance }}</p>
            <p>Valid Until: {{ validUntil }}</p>
        </div>
        <div class="user-options">
            <div class="sorting-options">
                <label>
                    Sort by:
                    <select v-model="sortOption" @change="sortTravelOptions">
                        <option value="duration">Duration</option>
                        <option value="price">Price</option>
                    </select>
                </label>
            </div>
            <div class="filtering-options">
                <label>
                    Filter by:
                </label>
                <custom-dropdown :options="uniqueCompanies" @input="filterTravelOptions"></custom-dropdown>

            </div>
        </div>
        <div class="table-div">
            <table>
                <thead>
                    <tr>
                        <th>Providers</th>
                        <th>Start Time</th>
                        <th>Duration</th>
                        <th>End Time</th>
                        <th>Price</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(option, index) in filteredTravelOptions" :key="index">
                        <td class="provider-data" @mouseover="showTooltip(index, $event, `provider`)"
                            @mouseout="hideTooltip">
                            {{
                                Array.from(new Set(option.providers.map(provider => provider.companyName))).join(', ')
                            }}
                        </td>
                        <td class="time-data">{{ option.FlightStart }}</td>
                        <td class="duration-data" @mouseover="showTooltip(index, $event, `duration`)"
                            @mouseout="hideTooltip">
                            {{ option.totalDuration }}
                        </td>
                        <td class="time-data">{{ option.FlightEnd }}</td>
                        <td class="price-data">{{ option.totalPrice }}</td>
                        <button @click="openBookingModal(index)" ref="bookingButton">Book</button>
                    </tr>
                </tbody>
            </table>
            <Tooltip :text="fullText" :showTooltip="tooltipVisible && hoveredIndex === hoveredIndex"
                :position="tooltipPosition" />
        </div>
    </div>
    <booking-modal :is-open="isBookingModalOpen" :booking-details="bookingDetails" @confirm="confirmBooking"
        @cancel="cancelBooking" ref="bookingModal"></booking-modal>
</template>

<script>
import SearchForm from '../components/SearchForm.vue';
import Tooltip from '../components/ToolTip.vue';
import CustomDropdown from '../components/CustomDropdown.vue';
import BookingModal from '../components/BookingModal.vue';

export default {
    components: {
        SearchForm,
        Tooltip,
        CustomDropdown,
        BookingModal,
    },
    data() {
        return {
            from: "",
            destination: "",
            travelOptions: [],
            distance: 0,
            validUntil: null,
            hoveredIndex: null,
            tooltipVisible: false,
            tooltipPosition: {},
            fullText: '',
            sortOption: 'duration',
            selectedCompanies: [],
            isBookingModalOpen: false,
            pricelistID: '',

            bookingDetails: {
                companyNames: [],
                startTime: '',
                firstName: '',
                lastName: '',
                totalPrice: 0,
                totalDuration: '',
                pricelistID: '',
                routes: [],
                from: '',
                destination: '',
            },
        };
    },
    beforeMount() {
        this.from = this.$route.query.from;
        this.destination = this.$route.query.destination;
        this.fetchFlights();
    },
    computed: {
        uniqueCompanies() {
            const companies = new Set();
            this.travelOptions.forEach(option => {
                option.providers.forEach(provider => {
                    companies.add(provider.companyName);
                });
            });
            return Array.from(companies);
        },
        filteredTravelOptions() {
            if (this.selectedCompanies.length === 0) {
                return this.travelOptions;
            } else {
                return this.travelOptions.filter(option => {
                    const companyNames = option.providers.map(provider => provider.companyName);
                    return this.selectedCompanies.some(company => companyNames.includes(company));
                });
            }
        },
    },
    methods: {
        async fetchFlights() {
            console.log("fetching");
            const response = await fetch(`http://localhost:8080/api/get/${this.from}/${this.destination}`);

            if (!response.ok) {
                this.$router.push({ name: "routeNotFound" });
                return;
            }

            const data = await response.json();
            console.log(data);

            const formattedProviders = data.possibleRoutes.map(leg => {
                return {
                    ...leg,
                    FlightStart: this.formatDate(leg.providers[0].flightStart),
                    FlightEnd: this.formatDate(leg.providers[leg.providers.length - 1].flightEnd),
                };
            });
            this.travelOptions = [...this.travelOptions, ...formattedProviders];
            this.distance = data.totalDistance;
            this.validUntil = this.formatDate(data.validUntil);
            this.pricelistID = data.pricelistID;
            this.sortTravelOptions();
        },

        formatDate(dateTimeString) {
            const parts = dateTimeString.split(/[-T:.Z]/);
            const date = new Date(
                Date.UTC(parts[0], parts[1] - 1, parts[2], parts[3], parts[4], parts[5])
            );
            return date.toLocaleString("en-GB"); // Adjust the locale as needed
        },
        showTooltip(index, event, type) {
            this.hoveredIndex = index;
            if (type === 'duration') {
                this.fullText = `${this.filteredTravelOptions[index].totalDuration}`;
            } else {
                this.fullText = Array.from(new Set(this.filteredTravelOptions[index].providers.map(provider => provider.companyName))).join(', ');
            }
            this.tooltipVisible = true;
            this.updateTooltipPosition(event);
            window.addEventListener('mousemove', this.updateTooltipPosition);
        },
        hideTooltip() {
            this.tooltipVisible = false;
            this.hoveredIndex = null;
            this.fullText = '';
            // Remove event listener
            window.removeEventListener('mousemove', this.updateTooltipPosition);
        },
        updateTooltipPosition(event) {
            const hoveredElement = event.target;
            const hoveredElementRect = hoveredElement.getBoundingClientRect();

            const topOffset = hoveredElementRect.top;

            this.tooltipPosition = {
                x: event.pageX + 10,
                y: topOffset + 10,
            };
        },
        sortTravelOptions() {
            if (this.sortOption === 'duration') {
                this.travelOptions.sort((a, b) => this.compareDurations(a.totalDuration, b.totalDuration));
            } else if (this.sortOption === 'price') {
                this.travelOptions.sort((a, b) => a.totalPrice - b.totalPrice);
            }
        },
        compareDurations(durationA, durationB) {
            const parseDuration = (duration) => {
                const parts = duration.split(', ')
                    .map(part => part.split(' '))
                    .map(([value, unit]) => ({ value: parseInt(value), unit }));

                return parts.reduce((acc, { value, unit }) => {
                    if (unit === 'days') {
                        acc += value * 24 * 60; // Convert days to minutes
                    } else if (unit === 'hours') {
                        acc += value * 60; // Convert hours to minutes
                    } else if (unit === 'minutes') {
                        acc += value;
                    }
                    return acc;
                }, 0);
            };

            const durationAInMinutes = parseDuration(durationA);
            const durationBInMinutes = parseDuration(durationB);

            return durationAInMinutes - durationBInMinutes;
        },
        filterTravelOptions(selectedCompanies) {
            this.selectedCompanies = Array.from(selectedCompanies)
            this.selectedCompanies.shift();
        },
        openBookingModal(index) {
            const selectedOption = this.filteredTravelOptions[index];

            this.bookingDetails = {
                companyNames: selectedOption.providers.map(provider => provider.companyName),
                startTime: selectedOption.providers[0].flightStart,
                firstName: '',
                lastName: '',
                totalPrice: parseFloat(selectedOption.totalPrice),
                totalDuration: selectedOption.totalDuration,
                pricelistID: this.pricelistID,
                routes: {
                    from: this.from,
                    destination: this.destination,
                },
                validUntil: this.validUntil,
            };
            this.isBookingModalOpen = true;
        },
        handleDocumentClick(event) {
            const bookingModalElement = this.$refs['bookingModal'].$el;
            const bookingButtonElem = this.$refs['bookingButton'];
            if (bookingButtonElem.some(button => button.contains(event.target))) {
                return;
            }
            if (bookingModalElement && !bookingModalElement.contains(event.target)) {
                this.cancelBooking();
            }
        },

        confirmBooking(bookingDetails) {
            console.log('Booking Confirmed:', bookingDetails);
            this.isBookingModalOpen = false;
        },
        cancelBooking() {
            this.isBookingModalOpen = false;
        }
    },
    watch: {
        isBookingModalOpen(newValue) {
            if (newValue) {
                this.$nextTick(() => {
                    window.addEventListener('click', this.handleDocumentClick);
                });
            } else {
                window.removeEventListener('click', this.handleDocumentClick);
            }
        },
    },
};
</script>

<style>
table {
    outline: 1px solid black;
    outline-offset: 2px;
    width: 100%;
    max-width: 1000px;
    table-layout: fixed;
}

td {
    border: 1px solid black;
    padding: 3px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    position: relative;
}

th {
    border: 1px solid black;
    padding: 3px;
}

.search-form-area {
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: rgb(28, 48, 79);
}

.info {
    display: flex;
    flex-direction: row;
    justify-content: center;
}

.info p {
    padding: 10px;
}

.big-div {
    justify-content: center;
    margin-bottom: 20px;
}

.table-div {
    display: flex;
    justify-content: center;
}

.sorting-options {
    display: flex;
    justify-content: center;
    margin-bottom: 20px;
}

.filtering-options {
    display: flex;
    justify-content: center;
    margin-bottom: 20px;
    align-items: center;
}

.user-options {
    display: flex;
    justify-content: center;
}
</style>
