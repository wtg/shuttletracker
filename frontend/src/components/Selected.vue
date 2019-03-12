<template>
	<div id='select'></div>
</template>

<script>
const Tabulator = require('tabulator-tables');
export default {
  data() {
    return {
      tabulator: null,
      tableColumn: [
        {title: 'Sunday', field: 'su', align: 'center', cellClick: function(e, cell){console.log(cell.getValue()); console.log(cell.getColumn().getField())},}, 
      {title: 'Monday', field: 'm', align: 'center', cellClick: function(e, cell){console.log(cell.getValue()); console.log(cell.getColumn().getField())},}, 
      {title: 'Tuesday', field: 't', align: 'center', cellClick: function(e, cell){console.log(cell.getValue()); console.log(cell.getColumn().getField())},}, 
      {title: 'Wednesday', field: 'w', align: 'center', cellClick: function(e, cell){console.log(cell.getValue()); console.log(cell.getColumn().getField())},}, 
      {title: 'Thursday', field: 'th', align: 'center', cellClick: function(e, cell){console.log(cell.getValue()); console.log(cell.getColumn().getField())},}, 
      {title: 'Friday', field: 'f', align: 'center', cellClick: function(e, cell){console.log(cell.getValue()); console.log(cell.getColumn().getField())},}, 
      {title: 'Saturday', field: 'sa', align: 'center', cellClick: function(e, cell){console.log(cell.getValue()); console.log(cell.getColumn().getField())},} 
      ],
    };
  },

  mounted() {
    this.tabulator = new Tabulator('#select', {
      data: [],
      columns: this.tableColumn,
      height: 250,
      layout: 'fitColumns',
      placeholder: 'Drag Times',
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
};
</script>

<style>
@import '~tabulator-tables/dist/css/tabulator.min.css'
</style>