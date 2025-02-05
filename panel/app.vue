<template>
  <v-card
    color="grey-lighten-4"
    rounded="0"
    flat
  >
    <v-toolbar collapse>
      <v-btn icon @click="refresh">
        <v-icon>mdi-refresh</v-icon>
      </v-btn>
      <v-btn icon @click="clear">
        <v-icon>mdi-delete</v-icon>
      </v-btn>
    </v-toolbar>
  </v-card>
  <div>
    <div id="timeline-container" ref="timelineContainer" style="height: 60dvh;">
    </div>
    <json-tree-view class="scrollable" style="height: 30dvh;" :json="traceValue">
    </json-tree-view>
  </div>
</template>

<script setup>
import { createVuetify } from 'vuetify'
import { Timeline, DataSet } from "vis-timeline/standalone"
import { onMounted } from "vue"
import { JsonTreeView } from "json-tree-view-vue3";
import 'json-tree-view-vue3/dist/style.css'

var timeline = ref(null)
const timelineContainer = ref(null)
const traceValue = ref("{}")

// Create a DataSet (allows two way data-binding)
var items = new DataSet([]);
var groups = new DataSet([
{
  id: 10,
  title: "traces",
  content: "traces",
  treeLevel: 1,
}
]);

// Configuration for the Timeline
var options = {
  minHeight: "100%",
  maxHeight: "100%",
  margin: { item: 0 },
  order: customOrder,
};

onMounted(() => {
  var container = document.getElementById("timeline-container")
  timeline = new Timeline(container, items, groups, options)
  timeline.on('select', onSelect)
})

function refresh() {
  $fetch('http://localhost:8866/api/get_traces').then(data => {
    console.debug(data)
    items.clear()
    var end = 0
    var start = 4000000000000
    data.forEach(trace => {
      const item = {}
      item.name = trace.name
      item.content = trace.name
      item.start = trace.start / 1000000
      item.end = item.start + trace.duration / 1000000
      item.group = 10
      item.trace = trace
      item.level = item.trace.level
      items.add(item)
      if (item.start < start) {
        start = item.start
      }
      if (item.end > end) {
        end = item.end
      }
    });
    console.debug(data, start, end)
    timeline.setOptions({"start": start, "end": end})
  })
}

async function clear() {
  await $fetch('http://localhost:8866/api/clear', { method: 'POST'})
  items.clear()
}

function onSelect(properties) {
  console.debug(properties)
  if (properties.items.length == 0) {
    traceValue.value = "{}"
    return
  }
  const item = items.get(properties?.items[0])
  traceValue.value = JSON.stringify(item?.trace)
}

function customOrder(a, b) {
  // order by id
  return a.level - b.level
}

</script>

<style scoped>
html{
    height: 100%;
}

body{
    margin:0;
    padding:0;
    overflow:hidden;
    height:100%;
}

#container{
    width:1000px;
    margin:0 auto;
    height:100%;
}
.scrollable {
   overflow-y: scroll;
}
</style>