<template>
  <NuxtLayout>
    <div>
      <!-- Static sidebar for desktop -->
      <div class="hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-72 lg:flex-col">
        <!-- Sidebar component, swap this element with another sidebar if you like -->
        <div class="flex grow flex-col gap-y-5 overflow-y-auto border-r border-gray-200 bg-white px-6 pb-4">
          <div class="flex h-16 shrink-0 items-center flex justify-center">
            <!--<img class="h-8 w-auto" src="/logo.svg" alt="" />-->
            <h1 class="text-2xl">DuckDB Studio</h1>
          </div>
          <ul>
            <li class="py-2">
              <TemplrJSDropdownList v-if="schemas.length > 0" :data="schemas" v-on:change="handleSelected_Schema" :show_label="true" name="schema" label="Schema" :selectedValue="schemas[2].name" />
            </li>
            <li class="py-2">
              <label for="table-search" class="sr-only">Search tables</label>
              <div class="mt-1 relative shadow-sm">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <!-- Search icon -->
                </div>
                <input type="text" name="table-search" id="table-search" class="h-8 block w-full border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-primary-600 sm:text-sm sm:leading-6" placeholder="Search tables..." />
              </div>
            </li>
          </ul>
          <nav class="flex flex-1 flex-col" style="height: calc(100vh - 100px); overflow-y: auto">
            <ul class="space-y-2">
              <li v-if="tables.length > 0">
                <h5 class="text-gray-600 dark:text-gray-400 uppercase font-bold px-3 mb-3">Tables</h5>
              </li>
              <li v-for="table in tables" :key="table">
                <a href="#" @click.prevent="get_table_data(table.name)" class="flex items-center text-base font-normal text-gray-900 dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                  <Icon name="icon-park-twotone:layout-one" class="h-4 w-4 text-center mr-2" />
                  <span class="text-sm truncate">{{ table.name }}</span>
                </a>
              </li>
            </ul>
          </nav>
        </div>

        <div class="flex-shrink-0 flex border-t border-r border-gray-300 p-4">
          <FooterMiniNav api_end_point="/footer-support-nav.json" />
        </div>
      </div>

      <div class="lg:pl-72">
        <div class="sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b border-gray-200 bg-white px-4 shadow-sm sm:gap-x-6 sm:px-6 lg:px-8">
          <!-- Separator -->
          <div class="w-full flex justify-end">
            <button @click="useNuxtApp().$bus.$emit('evtSideBarGuide')" class="h-8 rounded-2xl bg-primary-500 hover:bg-primary-700 text-white font-semibold px-4 border border-primary-500 hover:border-transparent" title="Practice SQL Queries">Practice SQL Queries</button>
          </div>

          <div class="h-6 w-px bg-gray-200 lg:hidden" aria-hidden="true"></div>
        </div>

        <main class="py-2">
          <div class="pl-2 pr-2">
            <!-- Your content -->
            <div class="flex-grow box-border border-2 border-gray-200">
              <!-- Code Editor -->
              <client-only placeholder="Codemirror Loading...">
                <codemirror :value="code" @ready="handleReady" @change="handleCodeChange" @focus="handleCodeFocus" @blur="handleCodeBlur" :placeholder="initialPaceholderSql" :style="{ height: auto, width: '100%' }" :autofocus="true" :indent-with-tab="true" :tab-size="2" :extensions="extensions" />
              </client-only>

              <!-- Action Buttons -->
              <footer class="flex justify-end gap-x-2 relative">
                <button @click="execute" class="absolute bottom-0 left-[45%] rounded-full bg-primary-100 hover:bg-primary-200 p-2 text-white z-40 h-10 w-10 transform -translate-x-1/2 translate-y-1/2 p-2 rounded-full z-10" title="Execute SQL code">
                  <span class="absolute inset-0 flex h-full w-full items-center justify-center" aria-hidden="true">
                    <Icon name="icon-park-twotone:play" class="w-10 h-10 text-primary-600" aria-label="Execute" />
                  </span>
                </button>
              </footer>
            </div>
            <div class="py-4">
              <footer class="flex justify-end gap-x-2 py-2">
                <button @click="clear" class="h-8 px-4 bg-transparent hover:bg-primary-500 text-primary-700 font-semibold hover:text-white border border-primary-500 hover:border-transparent" title="Reset SQL query editor">Clear</button>
                <button @click="rollback" class="h-8 bg-transparent hover:bg-primary-500 text-primary-700 font-semibold hover:text-white px-4 border border-primary-500 hover:border-transparent" title="Rollback to last SQL statement">Rollback</button>
                <button @click="prettify" class="h-8 bg-primary-500 hover:bg-primary-700 text-white font-semibold px-4 border border-primary-500 hover:border-transparent" title="Prettify SQL query using SQLGlot Python library">Prettify</button>
                <button @click="optimize" class="h-8 bg-primary-500 hover:bg-primary-700 text-white font-semibold px-4 border border-primary-500 hover:border-transparent" title="Optimize SQL query using SQLGlot Python library">Optimize</button>
                <!--Transpile pull-down menu button (start)-->
                <div class="relative inline-flex shadow-sm">
                  <button type="button" class="h-8 relative inline-flex items-center bg-primary-500 hover:bg-primary-700 text-white font-bold px-3 py-2 focus:z-10" @click="transpile(selectedDialect)" title="Transpile between SQL dialacts using SQLGlot Python library">Transpile</button>
                  <div class="relative -ml-px block">
                    <button type="button" class="h-8 relative inline-flex items-center bg-primary-500 hover:bg-primary-700 text-white font-bold px-2 py-2 focus:z-10" @click="isTranspileDropdownOpen = !isTranspileDropdownOpen">
                      <svg class="h-8 w-8" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
                      </svg>
                    </button>

                    <div v-if="isTranspileDropdownOpen" class="absolute right-0 z-10 mt-2 w-56 origin-top-right bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                      <div class="py-1" role="menu">
                        <a v-for="(dialect, index) in dialects" :key="index" @click="selectDialect(dialect)" class="text-gray-700 block px-4 py-2 text-sm cursor-pointer hover:bg-primary-100" role="menuitem">
                          {{ dialect }}
                        </a>
                      </div>
                    </div>
                  </div>
                </div>
                <!--Transpile pull-down menu button (end)-->
                <!--Parse pull-down menu button (start)-->
                <div class="relative inline-flex shadow-sm">
                  <button type="button" class="h-8 relative inline-flex items-center bg-primary-500 hover:bg-primary-700 text-white font-bold px-3 py-2 focus:z-10" @click="parsesql(selectedDialect)" title="Parse SQL and extract Table, Column and Projections using SQLGlot Python library">Parse SQL</button>
                  <div class="relative -ml-px block">
                    <button type="button" class="h-8 relative inline-flex items-center bg-primary-500 hover:bg-primary-700 text-white font-bold px-2 py-2 focus:z-10" @click="isParseSqlDropdownOpen = !isParseSqlDropdownOpen">
                      <svg class="h-8 w-8" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
                      </svg>
                    </button>

                    <div v-if="isParseSqlDropdownOpen" class="absolute right-0 z-10 mt-2 w-56 origin-top-right bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                      <div class="py-1" role="menu">
                        <a v-for="(metadata, index) in metadata_objects" :key="index" @click="selectParseSql(metadata)" class="text-gray-700 block px-4 py-2 text-sm cursor-pointer hover:bg-primary-100" role="menuitem">
                          {{ metadata }}
                        </a>
                      </div>
                    </div>
                  </div>
                </div>
                <!--Parse pull-down menu button (end)-->
              </footer>
            </div>
            <!-- Dynamic Data Table -->
            <div class="gap-y-3 p-4 lg:p-0">
              <div>
                <div class="mx-auto max-w-full pt-8 relative">
                  <!-- Progress bar (start)-->
                  <div class="absolute top-0 left-0 z-30 w-full h-full bg-gray-100 bg-opacity-75 transition-opacity" v-if="isLoading">
                    <div class="h-full bg-gray-100 bg-opacity-75 transition-opacity"></div>
                  </div>
                  <div class="animate-spin absolute top-1/2 left-1/2 z-40" v-if="isLoading">
                    <!-- SVG loader here -->
                    <svg class="h-8 w-8 text-zinc-800" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.324.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.24-.438.613-.431.992a6.759 6.759 0 010 .255c-.007.378.138.75.43.99l1.005.828c.424.35.534.954.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.57 6.57 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.28c-.09.543-.56.941-1.11.941h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.992a6.932 6.932 0 010-.255c.007-.378-.138-.75-.43-.99l-1.004-.828a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.214-1.281z" />
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                  </div>
                  <!-- Progress bar (end)-->
                  <div v-if="paginationLinks.length > 0" class="px-4">
                    <div class="sm:flex sm:items-center">
                      <div class="sm:flex-auto">
                        <h1 class="text-base font-semibold leading-6 text-gray-900">Query Results</h1>
                        <p class="mt-2 text-sm text-gray-700">DuckDB SQL query execution results</p>
                      </div>
                      <div class="mt-4 sm:ml-16 sm:mt-0 sm:flex-none">
                        <!--<button type="button" class="block  bg-primary-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600">Export</button>-->
                      </div>
                    </div>
                    <div class="mt-8 flow-root">
                      <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8 px-2">
                        <div class="inline-block min-w-full align-middle sm:px-6 max-h-96">
                          <table class="min-w-full divide-y divide-gray-300">
                            <thead class="bg-white border-b sticky top-0">
                              <tr>
                                <th scope="col" v-for="column in tableColumns" :key="column.field" class="sticky top-0 z-10 border-b border-gray-300 py-3.5 whitespace-nowrap px-4 text-left text-sm font-semibold text-zinc-600 sm:pl-0">
                                  {{ column.field }}
                                </th>
                              </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-200 bg-white h-96 overflow-y-auto">
                              <tr v-for="item in currentPageData" :key="item.id">
                                <td v-for="column in tableColumns" :key="column.field" class="whitespace-nowrap text-left text-sm text-zinc-600">
                                  {{ item[column.field] }}
                                </td>
                                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0"></td>
                              </tr>
                            </tbody>
                          </table>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <!-- Pagination -->
              <div v-if="paginationLinks.length > 0" class="mt-4 hidden sm:flex sm:flex-1 sm:items-center sm:justify-between px-4">
                <div>
                  <p class="text-sm text-gray-700">
                    Showing
                    <span class="font-medium">{{ $s.numberFormat((currentPage - 1) * itemsPerPage + 1, 0, '.', ',') }}</span>
                    to
                    <span class="font-medium">{{ $s.numberFormat(Math.min(currentPage * itemsPerPage, tableData.total_rows), 0, '.', ',') }}</span>
                    of
                    <span class="font-medium">{{ $s.numberFormat(tableData.total_rows, 0, '.', ',') }}</span>
                  </p>
                </div>

                <div>
                  <!-- Pagination Links -->
                  <nav class="inline-flex -space-x-px shadow-sm" aria-label="Pagination" id="nav_pagination">
                    <a href="#" @click.prevent="goToPage(currentPage - 1)" id="prev_page" class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                      <span class="sr-only">Previous</span>
                      <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd" />
                      </svg>
                    </a>

                    <!-- ... Loop through paginationLinks ... -->
                    <template v-for="pageLink in paginationLinks">
                      <a v-if="pageLink <= totalPages" :key="pageLink" href="#" :class="pageLinkClasses(pageLink)" :id="`pg_${pageLink}`" @click.prevent="goToPage(pageLink)">
                        {{ pageLink }}
                      </a>
                    </template>

                    <a href="#" @click.prevent="goToPage(currentPage + 1)" id="next_page" class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                      <span class="sr-only">Next</span>
                      <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd" />
                      </svg>
                    </a>
                  </nav>
                </div>
              </div>
              <div>&nbsp;</div>
              <!--Parse SQL resutls (start)-->

              <div class="flex flex-wrap gap-2" v-if="parsedSQLResults.length > 0">
                <span v-for="(item, index) in parsedSQLResults" :key="index" :class="`${getRandomColorClass()} text-gray-800 rounded-full px-4 py-1 text-sm font-medium`">
                  {{ item }}
                </span>
              </div>
              <!--Parse SQL results (end)-->
            </div>
          </div>
        </main>
      </div>
    </div>
    <AboutTemplrJS />
    <SideBarGuide form_title="DuckDB SQL Guide" form_description="Handy SQL queries to practice against the sample Tickit database" content_identifier="mastering-sql-queries-with-the-tickit-database-on-duckdb" class="z-40" />
  </NuxtLayout>
</template>
<script setup>
import { Codemirror } from 'vue-codemirror';
import { sql } from '@codemirror/lang-sql';
import { oneDark } from '@codemirror/theme-one-dark';
import _ from 'lodash';

import { computed } from 'vue';
// Function to get a random Tailwind CSS background color class
const getRandomColorClass = () => {
  const colorClasses = ['bg-red-200', 'bg-green-200', 'bg-blue-200', 'bg-yellow-200', 'bg-purple-200', 'bg-pink-200', 'bg-primary-200', 'bg-teal-200'];
  return colorClasses[Math.floor(Math.random() * colorClasses.length)];
};
//-----------------------------------------------------------------------------------------------------
//1. Page level variable declarations
//-----------------------------------------------------------------------------------------------------
const data = ref({});
const schemas = ref([]);
const selectedSchema = ref('');
const tables = ref([]);
const parsedSQLResults = ref([]);
const initialPaceholderSql = `
SELECT venues.venuename, categories.catname, events.eventname,listings.listtime, (sales.saletime- listings.listtime) as sale_age
FROM sales,listings, venues,events,categories
WHERE
events.eventid=sales.eventid
and
events.catid=categories.catid
and sales.listid=listings.listid
and venues.venueid=events.venueid
ORDER BY sale_age desc, sales.saletime desc,  listings.listtime desc, venues.venuename,  categories.catname,  events.eventname
LIMIT 20
`;

//-----------------------------------------------------------------------------------------------------
//2. Page level events handlers
//-----------------------------------------------------------------------------------------------------

function handleSelected_Schema(e) {
  getTables(e.target.value);
  selectedSchema.value = e.target.value;
  console.log(`Select Value = [${e.target.value}]`);
  reInitializeDataTable();
}

//-----------------------------------------------------------------------------------------------------
//3. Async api call to get schemas from duckdb and set the result to bind the dropdownlist
//-----------------------------------------------------------------------------------------------------
try {
  const schemaRequestOptions = {
    method: 'POST',
    body: {
      query: 'select name, name as code from (select distinct(table_schema) as name from information_schema.tables order by table_schema asc)',
    },
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const schema_response = await $fetch(`${useRuntimeConfig().public.DUCKDB_DATA_API_BASE_PATH}/execute/sql`, schemaRequestOptions);

  if (schema_response.data && Array.isArray(schema_response.data)) {
    schemas.value = schema_response.data;
  }
} catch (e) {
  console.error(e);
}

//-----------------------------------------------------------------------------------------------------
//4. Functions
//-----------------------------------------------------------------------------------------------------

const getTables = async (schema) => {
  try {
    const requestOptions = {
      method: 'POST',
      body: {
        query: `SELECT table_name as name FROM information_schema.tables where table_schema='${schema}' order by table_name asc`,
      },
      headers: {
        'Content-Type': 'application/json',
      },
    };
    const response = await $fetch(`${useRuntimeConfig().public.DUCKDB_DATA_API_BASE_PATH}/execute/sql`, requestOptions);

    if (response.data && Array.isArray(response.data)) {
      tables.value = response.data;
    }
  } catch (e) {
    console.error(e);
  }
};

// Codemirror EditorView instance ref
const extensions = [sql()];
const view = shallowRef();
const isExecuting = ref(false); // Reactive variable to track execution state
const isDataTableVisible = ref(false); // Default to false
const isLoading = ref(false);

// Status is available at all times via Codemirror EditorView
const getCodemirrorStates = () => {
  const state = view.value.state;
  const ranges = state.selection.ranges;
  const selected = ranges.reduce((r, range) => r + range.to - range.from, 0);
  const cursor = ranges[0].anchor;
  const length = state.doc.length;
  const lines = state.doc.lines;
  // more state info ...
  // return ...
};

const handleReady = (payload) => {
  view.value = payload.view;
};

const handleCodeChange = (newValue) => {
  // Handle code change event here
  // For example, you might want to update the 'code' data property
  code.value = newValue;
};

const handleCodeFocus = (event) => {
  // Handle code focus event here
};

const handleCodeBlur = (event) => {
  // Handle code blur event here
};

const clear = () => {
  parsedSQLResults.value = [];
  // Clear the code in the codemirror editor
  reInitializeDataTable(); // Reset the data table
  if (view.value) {
    //console.log(view.value.state.doc.length);
    view.value.dispatch({
      changes: { from: 0, to: view.value.state.doc.length, insert: '' },
    });
    //console.log(view.value.state.doc.length);
  }
};

const get_table_data = (table_name) => {
  const sql = `select * from ${selectedSchema.value}.${table_name} limit 100;`;
  console.log(sql);
  if (view.value) {
    view.value.dispatch({
      changes: { from: 0, to: view.value.state.doc.length, insert: sql },
    });
    //console.log(view.value.state.doc.length);
    //console.log(view.value.state.doc.toString());
    execute();
  }
};

const execute = () => {
  parsedSQLResults.value = [];
  // Execute the code from the codemirror editor
  const codeToExecute = view.value ? view.value.state.doc.toString() : '';
  if (view.value.state.doc.length === 0) {
    useNuxtApp().$toast.show({
      type: 'danger',
      message: `Empty SQL query is not allowed.`,
      position: 'top-left',
      timeout: 2000,
    });
    return;
  }

  reInitializeDataTable(); // Reset the data table
  // Call fetchDataForPage(1) to trigger data fetching based on the code
  fetchDataForPage(1);
};

const rollback = () => {
  reInitializeDataTable(); // Reset the data table
  if (view.value) {
    view.value.dispatch({
      changes: { from: 0, to: view.value.state.doc.length, insert: code.value },
    });
  }
};

const reInitializeDataTable = () => {
  currentPageData.value = [];
  currentPage = 1; // Initialize currentPage to 1
  totalPages.value = 0;
  tableData.value = {};
  tableColumns.value = [];
  paginationLinks.value = [];
};

const isCodeEmpty = computed(() => {
  return !view.value || view.value.state.doc.length === 0;
});

// Initialize your 'code' data property here
const code = ref(''); // Initialize with your default code

//--------------------------------------------------------------------------Data Table

// Initialize current page data with the first page

const currentPageData = ref([]);

let currentPage = 1; // Initialize currentPage to 1
const totalPages = ref(0);
const tableData = ref({});
const tableColumns = ref([]);
const paginationLinks = ref([]);

const fetchDataForPage = async (page) => {
  isLoading.value = true;
  // Introduce an artificial delay for debugging (2 seconds)
  //await new Promise((resolve) => setTimeout(resolve, 2000));

  // Make sure the code is not empty before fetching data
  if (view.value.state.doc.length > 0) {
    // Construct the dynamic API endpoint based on the code entered
    //const dynamicEndpoint = `http://localhost:8080/${code.value}?limit=${itemsPerPage}&offset=${(page - 1) * itemsPerPage}`;

    try {
      isDataTableVisible.value = true;
      const requestOptions = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: { query: code.value },
      };

      //const response = await $fetch(dynamicEndpoint);

      const response = await $fetch(`${useRuntimeConfig().public.DUCKDB_DATA_API_BASE_PATH}/execute/sql`, requestOptions);
      //const responseData = await response.json();
      //console.log('response=', JSON.stringify(response));
      if (response) {
        if (response.data && Array.isArray(response.data)) {
          const startIndex = (page - 1) * itemsPerPage;
          const endIndex = startIndex + itemsPerPage;
          currentPageData.value = response.data.slice(startIndex, endIndex);
          tableData.value = response;
          currentPage = page; // Update currentPage when a page is clicked
          const totalRows = response.total_rows;
          totalPages.value = Math.ceil(totalRows / itemsPerPage); // Corrected calculation
          console.log(response.data[0]);
          // Build dynamic column definitions from JSON response attributes
          tableColumns.value = Object.keys(response.data[0]).map((attribute) => {
            return {
              label: useNuxtApp().$s.humanize(attribute),
              field: attribute,
            };
          });

          // Calculate pagination links
          paginationLinks.value = Array.from({ length: totalPages.value }, (_, index) => index + 1);

          //console.log('paginationLinks=', JSON.stringify(paginationLinks.value));

          updatePaginationLinks(currentPage, totalRows); // Pass totalRows instead of response.data.total_rows
          isLoading.value = false;
        } else {
          console.error('Inside error');
          useNuxtApp().$toast.show({
            type: 'danger',
            message: `${response.error}`,
            position: 'top-left',
            timeout: 2000,
          });
          isLoading.value = false;
          reInitializeDataTable(); // Reset the data table
        }
      } else {
        useNuxtApp().$toast.show({
          type: 'danger',
          message: 'No response from server',
          position: 'top-left',
          timeout: 2000,
        });
        isLoading.value = false;
        reInitializeDataTable(); // Reset the data table
      }
    } catch (error) {
      console.error(error);
      useNuxtApp().$toast.show({
        type: 'danger',
        error: error,
        message: `${error}`,
        position: 'top-left',
        timeout: 2000,
      });
      isLoading.value = false;
      reInitializeDataTable(); // Reset the data table
    }
  }
};

const itemsPerPage = 25; // Set the number of items per page
const maxPaginationLinks = 10; // Set the maximum number of pagination links

const updatePaginationLinks = (currentPage, totalRows) => {
  //console.log('updatePaginationLinks - currentPage:', currentPage);
  //console.log('updatePaginationLinks - totalRows:', totalRows);

  const totalPageCount = Math.ceil(totalRows / itemsPerPage);
  //console.log('updatePaginationLinks - totalPageCount:', totalPageCount);

  // Calculate middlePage and adjust it if it exceeds totalPageCount
  const pageLinksToShow = Math.min(totalPageCount, maxPaginationLinks);
  const middlePage = Math.floor(pageLinksToShow / 2);
  //console.log('updatePaginationLinks - middlePage:', middlePage);

  // Calculate startPage and endPage considering middlePage and totalPageCount
  let startPage = Math.max(1, currentPage - middlePage);
  let endPage = Math.min(totalPageCount, currentPage + middlePage);

  // Adjust startPage and endPage to ensure they stay within bounds
  if (endPage - startPage + 1 > pageLinksToShow) {
    if (currentPage <= middlePage) {
      endPage = startPage + pageLinksToShow - 1;
    } else if (totalPageCount - currentPage < middlePage) {
      startPage = totalPageCount - pageLinksToShow + 1;
    } else {
      startPage = currentPage - Math.floor(pageLinksToShow / 2);
      endPage = startPage + pageLinksToShow - 1;
    }
  }

  paginationLinks.value = Array.from({ length: endPage - startPage + 1 }, (_, index) => startPage + index);

  //console.log('updatePaginationLinks - paginationLinks:', paginationLinks.value);
};

// Function to determine pagination link classes
const pageLinkClasses = (pageLink) => {
  const baseClasses = 'relative inline-flex items-center px-4 py-2 text-sm font-semibold ring-1 ring-inset ring-gray-300 focus:outline-offset-0';
  if (pageLink === currentPage) {
    return `${baseClasses} bg-primary-600 text-white focus:z-20 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600`;
  } else {
    return `${baseClasses} text-gray-900 hover:bg-gray-50 focus:z-20 focus:outline-offset-0`;
  }
};

const goToPage = async (page) => {
  if (page < 1 || page > paginationLinks.length || page > totalPages.value) {
    return; // Don't go out of bounds
  }
  //currentPage < totalPages

  fetchDataForPage(page);
};

const transpile = async (dialect) => {
  // Get the code from the CodeMirror editor
  const codeToTranspile = view.value ? view.value.state.doc.toString() : '';
  const requestOptions = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: { sql: codeToTranspile, transpile_to: dialect.toLowerCase() },
  };
  try {
    const response = await $fetch(`${useRuntimeConfig().public.DUCKDB_DATA_API_BASE_PATH}/sqlglot/transpile`, requestOptions);
    if (response) {
      if (view.value) {
        view.value.dispatch({
          changes: {
            from: 0,
            to: view.value.state.doc.length,
            insert: response.result_sql,
          },
        });
      }
    } else {
      useNuxtApp().$toast.show({
        type: 'danger',
        message: 'Error',
        position: 'top-left',
        timeout: 2000,
      });
    }
  } catch (error) {
    console.error('Error transpiling SQL:', error);
    // Handle the error appropriately
  }
};

const prettify = async () => {
  // Get the code from the CodeMirror editor
  const input = view.value ? view.value.state.doc.toString() : '';
  const requestOptions = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: { sql: input },
  };
  try {
    const response = await $fetch(`${useRuntimeConfig().public.DUCKDB_DATA_API_BASE_PATH}/sqlglot/prettify`, requestOptions);
    if (response) {
      if (view.value) {
        view.value.dispatch({
          changes: {
            from: 0,
            to: view.value.state.doc.length,
            insert: response.result_sql,
          },
        });
      }
    } else {
      useNuxtApp().$toast.show({
        type: 'danger',
        message: 'Error',
        position: 'top-left',
        timeout: 2000,
      });
    }
  } catch (error) {
    console.error('Error prettifying SQL:', error);
    // Handle the error appropriately
  }
};

const parsesql = async (metadata) => {
  //alert(metadata.replace('Extract ', '').replace('(s)', '').toLowerCase());

  parsedSQLResults.value = [];
  // Clear the code in the codemirror editor
  reInitializeDataTable(); // Reset the data table

  let operation = metadata.replace('Extract ', '').replace('(s)', '').toLowerCase();
  // Get the code from the CodeMirror editor
  const input = view.value ? view.value.state.doc.toString() : '';
  const requestOptions = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: { sql: input },
  };
  try {
    const response = await $fetch(`${useRuntimeConfig().public.DUCKDB_DATA_API_BASE_PATH}/sqlglot/extract/${operation}`, requestOptions);
    if (response) {
      parsedSQLResults.value = [...new Set(response.data)];
    } else {
      useNuxtApp().$toast.show({
        type: 'danger',
        message: 'Error',
        position: 'top-left',
        timeout: 2000,
      });
    }
  } catch (error) {
    console.error(`Error extracting ${operation}s`, error);
    // Handle the error appropriately
  }
};

const isTranspileDropdownOpen = ref(false);
const isParseSqlDropdownOpen = ref(false);
const selectedDialect = ref('Postgres'); // Default selected dialect
const selectedParseSql = ref('table'); // Default selected metadata
const dialects = ['athena', 'bigQuery', 'clickHouse', 'databricks', 'doris', 'drill', 'duckDB', 'hive', 'mySQL', 'oracle', 'postgres', 'presto', 'pRQL', 'redshift', 'snowflake', 'spark', 'spark2', 'sQLite', 'starRocks', 'tableau', 'teradata', 'trino', 'tSQL'];

const metadata_objects = ['Extract Table(s)', 'Extract Column(s)', 'Extract Projection(s)'];

const selectDialect = (dialect) => {
  selectedDialect.value = dialect;
  transpile(dialect); // Transpile immediately on selection
  isTranspileDropdownOpen.value = false; // Close the dropdown
};

const selectParseSql = (metadata) => {
  selectedParseSql.value = metadata;
  parsesql(metadata); // Transpile immediately on selection
  isParseSqlDropdownOpen.value = false; // Close the dropdown
};

const editorHeight = ref('200px'); // Start with a default height, e.g., 200px

onMounted(async () => {
  if (schemas.value.length > 0) {
    selectedSchema.value = schemas.value[0].name;
    getTables(schemas.value[0].name);
  }
});

const isOpen = ref(false);

useNuxtApp().$bus.$on('evtAboutTemplrJS', (data) => {
  //alert("Hello");
  isOpen.value = true;
});
</script>
<style>
.animate-spin {
  animation: spin 2s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
