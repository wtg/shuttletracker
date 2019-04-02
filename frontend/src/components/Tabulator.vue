<template>
	<div id='table'></div>
</template>

<script>
// npm install tabulator-tables --save
import Vue from 'vue';
import EventBus from '../event_bus.ts';
const Tabulator = require('tabulator-tables');
export default Vue.extend({
  data() {
    return {
      tabulator: null,
      tableColumn: [
      {title: 'Sunday', field: 'sunday', align: 'center', headerSort:false, cellClick: this.sendData}, 
      {title: 'Monday', field: 'monday', align: 'center', headerSort:false, cellClick: this.sendData}, 
      {title: 'Tuesday', field: 'tuesday', align: 'center', headerSort:false, cellClick: this.sendData}, 
      {title: 'Wednesday', field: 'wednesday', align: 'center', headerSort:false, cellClick: this.sendData}, 
      {title: 'Thursday', field: 'thursday', align: 'center', headerSort:false, cellClick: this.sendData}, 
      {title: 'Friday', field: 'friday', align: 'center', headerSort:false, cellClick: this.sendData}, 
      {title: 'Saturday', field: 'saturday', align: 'center', headerSort:false, cellClick: this.sendData} ],
    
    tableData: [ 
    {id: 12, sunday: '7:00', monday: '7:00', tuesday: '7:00', wednesday: '7:00', thursday: '7:00', friday: '7:00', saturday: '7:00'}, 
    {id: 1, sunday: '8:00', monday: '8:00', tuesday: '8:00', wednesday: '8:00', thursday: '8:00', friday: '8:00', saturday: '8:00'}, 
    {id: 2, sunday: '9:00', monday: '9:00', tuesday: '9:00', wednesday: '9:00', thursday: '9:00', friday: '9:00', saturday: '9:00'}, 
    {id: 3, sunday: '10:00', monday: '10:00', tuesday: '10:00', wednesday: '10:00', thursday: '10:00', friday: '10:00', saturday: '10:00'}, 
    {id: 4, sunday: '11:00', monday: '11:00', tuesday: '11:00', wednesday: '11:00', thursday: '11:00', friday: '11:00', saturday: '11:00'}, 
    {id: 5, sunday: '12:00', monday: '12:00', tuesday: '12:00', wednesday: '12:00', thursday: '12:00', friday: '12:00', saturday: '12:00'}, 
    {id: 6, sunday: '13:00', monday: '13:00', tuesday: '13:00', wednesday: '13:00', thursday: '13:00', friday: '13:00', saturday: '13:00'}, 
    {id: 7, sunday: '14:00', monday: '14:00', tuesday: '14:00', wednesday: '14:00', thursday: '14:00', friday: '14:00', saturday: '14:00'}, 
    {id: 8, sunday: '15:00', monday: '15:00', tuesday: '15:00', wednesday: '15:00', thursday: '15:00', friday: '15:00', saturday: '15:00'}, 
    {id: 9, sunday: '16:00', monday: '16:00', tuesday: '16:00', wednesday: '16:00', thursday: '16:00', friday: '16:00', saturday: '16:00'}, 
    {id: 10, sunday: '17:00', monday: '17:00', tuesday: '17:00', wednesday: '17:00', thursday: '17:00', friday: '17:00', saturday: '17:00'}, 
    {id: 11, sunday: '18:00', monday: '18:00', tuesday: '18:00', wednesday: '18:00', thursday: '18:00', friday: '18:00', saturday: '18:00'}, ],
	  addData: {id: null, sunday: null, monday: null, tuesday: null, wednesday: null, thursday: null, friday: null, saturday: null},
	  };
	},
  mounted() {
    this.tabulator = new Tabulator('#table', {
	  data: this.tableData,
	  columns: this.tableColumn,
	  layout: 'fitColumns',
	  });
	},
  methods: {
    addTabulator() {
	  const obj = Object.assign({}, this.addData);
	  obj.id = this.tableData.length;
	  this.tableData.push(obj);
	  Object.keys(this.addData).forEach((key => {
	    this.addData[key] = null;
      }));
    },

    sendData : function(e, cell) {
      const payload = {
        time: cell.getValue(),
        day: cell.getColumn().getField()
      }
      EventBus.$emit('TIME_SENT', payload);
    },
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
#table {
  background-color:#666;
  border: 1px solid #333;
  border-radius: 10px;
}

#table .tabulator-header {
  background-color:#666;
  color:#fff;
}

#table .tabulator-header .tabulator-col,
#table .tabulator-header .tabulator-col-row-handle {
  white-space: normal;
  background-color:#333;
}

#table .tabulator-tableHolder .tabulator-table .tabulator-row {
  background-color:#666;
  color:#fff;
}

#table .tabulator-tableHolder .tabulator-table .tabulator-row:nth-child(even) {
  background-color:#444;
}
</style>