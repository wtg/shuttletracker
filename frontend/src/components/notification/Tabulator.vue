<template>
	<div id='table'></div>
</template>

<script>
// TO BUILD TABULATOR NEED ===>
// ===>  "npm install tabulator-tables --save" <====
import Vue from 'vue';
import EventBus from '../../event_bus.ts';
const Tabulator = require('tabulator-tables');
export default Vue.extend({
  data() {
    return {
      days: ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'],
      times: [],
      tabulator: null,
      tableColumn: [],
      tableData: [],
      addData: {id: null, Sunday: null, Monday: null, Tuesday: null, Wednesday: null, Thursday: null, Friday: null, Saturday: null},
	  };
	},
  created() {
    //tableColumn
    var tc = [];
    for ( var i = 0; i < this.days.length; i++ ) {
      let obj = {title: this.days[i], field: this.days[i], align: 'center', headerSort:false, cellClick: this.sendData};
      tc.push(obj);
    }
    this.tableColumn = tc;
    //times
    tc = [];
    var start_hour = 7; //first hour
    var end_hour = 19; //last hour
    var minutes = ['00']; //minutes
    for ( var i = start_hour; i <= end_hour; i++ ) {
      for ( var j = 0; j < minutes.length; j++ ) {
        tc.push(i.toString() + ":" + minutes[j]);
      }
    }
    this.times = tc;
    //tableData
    tc = [];
    for ( var i = 0; i < this.times.length; i++ ) {
      let obj = {id: i, Sunday: this.times[i], Monday: this.times[i], Tuesday: this.times[i], Wednesday: this.times[i], Thursday: this.times[i], Friday: this.times[i], Saturday: this.times[i]};
      tc.push(obj);
    }
    this.tableData = tc;
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