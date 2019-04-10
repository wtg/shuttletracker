<template>
	<div id='select'></div>
</template>

<script>
import Vue from 'vue';
import EventBus from '../../event_bus.ts';
const Tabulator = require('tabulator-tables');
export default Vue.extend({
  data() {
    return {
      tabulator: null,
      tableColumn: [
        {title: 'Day', field: 'day', align: 'center', headerSort:false,cellClick: this.delete}, 
        {title: 'Time', field: 'time', align: 'center', headerSort:false,cellClick: this.delete}, 
      ],
      counter: 0,
      days: ['Sunday','Monday','Tuesday','Wednesday','Thurday','Friday','Saturday'],
      selected_times: [],
    };
  },
  methods: {
    receiveData (payload) {
      let selected_time = payload.day + ' ' + payload.time;
      if ( this.selected_times.indexOf( selected_time ) == -1 ) {
        this.counter += 1;
        this.tabulator.addData([{id:this.counter, day:payload.day.charAt(0).toUpperCase()+payload.day.slice(1), time:payload.time}], true);
        this.selected_times.push(selected_time);
        EventBus.$emit('TIME_ADDED', payload);
      }
    },
    delete : function(e, cell) {
      let selected_time = cell.getRow().getCells();
      const payload = {
        time: selected_time[1].getValue(),
        day: selected_time[0].getValue().charAt(0).toLowerCase() + selected_time[0].getValue().slice(1)
      }
      EventBus.$emit('TIME_REMOVED', payload);
      selected_time = payload.day + ' ' + payload.time;
      let index = this.selected_times.indexOf( selected_time );
      this.selected_times[index] = this.selected_times[0];
      this.selected_times.shift();
      cell.getRow().delete();
    },
    reset () {
      for (let row in this.tabulator.getRows()) {
        row.delete();
      }
    },
  },
  mounted() {
    this.tabulator = new Tabulator('#select', {
      data: [],
      columns: this.tableColumn,
      height: 250,
      layout: 'fitColumns',
      placeholder: 'Select Times',
    });
    EventBus.$on('TIME_SENT', (payload) => {
      this.receiveData(payload)
    });
  },
  watch: {
    tableData: {
      handler: function(newData) {
        this.tabulator.replaceData(newData);
      },
      deep: true,
    },
  },
});
</script>

<style>
@import '~tabulator-tables/dist/css/tabulator_midnight.min.css'
</style> 