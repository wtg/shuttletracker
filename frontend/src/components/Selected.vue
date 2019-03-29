<template>
	<div id='select'>
   <button v-on:click="reset()">Reset</button> 
  </div>
</template>

<script>
import Vue from 'vue';
import EventBus from '../event_bus.ts';
const Tabulator = require('tabulator-tables');
export default Vue.extend({
  data() {
    return {
      tabulator: null,
      tableColumn: [
        {title: 'Day', field: 'day', align: 'center', headerSort:false,cellClick: function(e, cell){
            cell.getRow().delete();
        }}, 
        {title: 'Time', field: 'time', align: 'center', headerSort:false,cellClick: function(e, cell){
           cell.getRow().delete();
        }}, 
      ],
      counter: 0,
      days: ['Sunday','Monday','Tuesday','Wednesday','Thurday','Friday','Saturday'],
    };
  },

  methods: {
    receiveData (payload) {
      this.counter += 1;
      this.tabulator.addData([{id:this.counter, day:payload.day.charAt(0).toUpperCase()+payload.day.slice(1), time:payload.time}], true);
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