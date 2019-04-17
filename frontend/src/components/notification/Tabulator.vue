<template>
	<div id='table'></div>
</template>

<script>
// TO BUILD TABULATOR NEED ===>
// ===>  "npm install tabulator-tables --save" <====
import Vue from 'vue';
import EventBus from '@/event_bus.ts';
const Tabulator = require('tabulator-tables');
export default Vue.extend({
  props: ['stop_id', 'route'],
  data() {
    return {
      times: {},
      tabulator: null,
      tableColumn: [],
      tableData: [],
      addData: {},
	  };
	},
  computed: {
    curr_route: function() {
      let cr;
      switch ( this.route ) {
        case '1':
          cr = 'W';
        case '2':
          cr = 'E';
        default:
          cr = 'N';
      }
      return cr;
    },
    days: function() {
      let d = [];
      var shuttleTimes = this.selectFile();
      Object.keys(shuttleTimes).forEach(function(k){
        if ( k !== 'stopname' ) {
          d.push(k);
          console.log(k);
        }
      });
      return d;
    },
  },
  created() {
    //accessing schedule
    var shuttleTimes = this.selectFile();
    
    //tableColumn
    var tc = [];
    for ( var i = 0; i < this.days.length; i++ ) {
      let name = this.days[i].split('_')
      let obj = {title: name[0]+' '+name[1], field: this.days[i], align: 'center', headerSort:false, cellClick: this.sendData};
      tc.push(obj);
    }
    this.tableColumn = tc;

    //weekday times
    var t = [];
    tc = {};
    var i = 0, maxLength = 0;
    const now = new Date();
    Object.keys(shuttleTimes).forEach(function(k){
      let Length = 0;
      if ( k !== 'stopname' ) {
        t = [];
        shuttleTimes[k].forEach(function(element) {
          if ( Math.abs( Number(element.split(':')[0]) - now.getHours() ) < 3 )
            t.push(element);
            Length += 1;
        });
        if ( maxLength < Length ) maxLength = Length;
        tc[k] = t;
        i++;
      }
    });
    this.times = tc;
    //console.log(this.times);

    //tableData
    tc = [];
    //create obj format
    for ( var i = 0; i < maxLength; i++ ) {
      let obj = {};
      for ( var i = 0; i < this.times.length; i++ ) {
        obj[this.days[i]] = null;
      }
      obj['id'] = i;
      Object.keys(this.times).forEach(function(k){
        obj[k] = this.times[k][i];
      });
      tc.push(obj)
    }
    /*for ( var i = 0; i < this.times.length; i++ ) {
      for ( var j = 0; j < this.times[i].length; j++ ) {

      }
      /*let obj = {id: i};
      for ( var j = 0; j < this.times[i].length; j++) {
        obj[this.days[]]
      }
      let obj = {id: i, Weekday: this.times[i], Weekend: this.times[i]};
      tc.push(obj);
    }*/
    this.tableData = tc;
  },
  mounted() {
    this.tabulator = new Tabulator('#table', {
	  data: this.tableData,
	  columns: this.tableColumn,
	  layout: 'fitColumns',
    placeholder: 'Select Route',
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
    selectFile() {
      let rt;
      switch ( this.route ) {
      case '1': rt = require('@/assets/shuttle_times/stops/1.json');
      default : rt = require('@/assets/shuttle_times/stops/1.json');
      //TODO add more schedules
      }
      return rt;
    },
    sendData : function(e, cell) {
      const payload = {
        time: cell.getValue(),
        day: cell.getColumn().getField()
      }
      EventBus.$emit('TIME_SENT', payload);
      console.log(this.curr_route);
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