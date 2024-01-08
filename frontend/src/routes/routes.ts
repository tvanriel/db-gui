import {type RouteRecordRaw} from 'vue-router';
import Homepage from '../components/Homepage.vue';
import AddConnection from '../components/AddConnection.vue';
import AddConnectionScript from '../components/AddConnectionScript.vue';
import Schema from '../components/Schema.vue';
import ListTables from '../components/ListTables.vue';
import SqlEditor from '../components/SqlEditor.vue';
import DescribeTable from '../components/DescribeTable.vue';
import Select from '../components/Select.vue';
import Import from '../components/Import.vue';
import Export from '../components/Export.vue';
import CreateTable from '../components/CreateTable.vue';

export const routes: RouteRecordRaw[] = [
	{path: '/', component: Homepage},
	{path: '/add-connection', component: AddConnection},
	{path: '/add-connection-script', component: AddConnectionScript},
	{path: '/schema', component: Schema},
	{path: '/list-tables', component: ListTables},
	{path: '/describe-table', component: DescribeTable},
	{path: '/select', component: Select},
	{path: '/sql-editor', component: SqlEditor},
	{path: '/import', component: Import},
	{path: '/export', component: Export},
	{path: '/create-table', component: CreateTable},
];
