<script lang="ts">
import { SelectionConfig } from '../services/select-service';
import { useAuthStore, useSchemaStore, useSelectStore } from '../store/store';


export default {
    data: () => ({
        fieldsetSelectOpen: false,
        fieldsetScopeOpen: false,
        fieldsetSearchOpen: false,
        fieldsetSortOpen: false,
        selectFunc: "",
        selectName: "",
        whereCol: "",
        whereOperator: "",
        whereVal: "",
    }),

    setup() {
        return {
            authStore: useAuthStore(),
            schemaStore: useSchemaStore(),
            selectStore: useSelectStore(),
        }
    },
    computed: {
        config(): SelectionConfig {
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            return this.selectStore.selectConfigFor(
                token, db.name, db.selectedTable,
            );
        },

        order() {
            return this.config.order;
        },

        table() {
            const token = this.authStore.current
            
            if (token === undefined) {
                return undefined;
            }
            
            const db = this.schemaStore.getCurrentDatabase(token);
            
            if (db === undefined) {
                return null;
            }

            return db.tables[db.selectedTable] ?? null;
        },

        columns() {
            return this.config.columns;
        },

        wheres() {
            return this.config.where;
        },

        columnFunc(i: number): string {
            return this.columns[i].func;
        },

        columnName(i: number): string {
            return this.columns[i].name;
        },

        limit() {
            return this.config.limit.limit;
        },

        offset() {
            return this.config.limit.offset;
        }
    },

    methods: {
        openFieldsetSelect() {
            this.fieldsetSelectOpen = !this.fieldsetSelectOpen
        },
        openFieldsetScope() {
            this.fieldsetScopeOpen = !this.fieldsetScopeOpen
        },

        openFieldsetSearch() {
            this.fieldsetSearchOpen = !this.fieldsetSearchOpen
        },
        openFieldsetSort() {
            this.fieldsetSortOpen = !this.fieldsetSortOpen
        },

        setColumnFunc(id: number, to: string) {
            const columns = this.columns;
            columns[id].func = to;
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.selectStore.updateCols(
                token, db.name, db.selectedTable,
                columns,    
            )    


        },

        addColumnFunc(func: string) {
            const columns = this.columns;

            columns.push({name: "", func: func});
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.selectStore.updateCols(
                token, db.name, db.selectedTable,
                columns,
            )
            this.selectFunc = ""
            this.selectName = ""
        },

        addColumnName(name: string) {
            const columns = this.columns;

            columns.push({name: name, func: ""});
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
this.selectStore.updateCols(
                token, db.name, db.selectedTable,
                columns,
            )
            this.selectFunc = ""
            this.selectName = ""
        },

        removeWhere(i: number) {
            const where = this.wheres;
            where.splice(i,1);
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.selectStore.updateWhere(
                token, db.name, db.selectedTable,
                where,    
            )
        },

        removeColumn(i: number) {
            const columns = this.columns;
            columns.splice(i,1);
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.selectStore.updateCols(
                token, db.name, db.selectedTable,
                columns,    
            )        
        },

        setColumnName(id: number, to: string) {
            const columns = this.columns;
            columns[id].name = to;
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.selectStore.updateCols(
                token, db.name, db.selectedTable,
                columns,    
            )    
        },

        setLimit(to: string) {
            const limit = this.config.limit;
            limit.limit = parseInt(to);
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.selectStore.updateLimit(
                token, db.name, db.selectedTable,
                limit,
            )

        },

        setOffset(to: string) {
            const limit = this.config.limit;
            limit.offset = parseInt(to);
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.selectStore.updateLimit(
                token, db.name, db.selectedTable,
                limit,    
            )         
        },

        select() {
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.selectStore.executeSelect(
                token, db.name, db.selectedTable,
            )         
            window.scrollTo(0,0);
        },

        addWhereCol(columnName: string) {

            const wheres = this.wheres;

            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            wheres.push({columnName: columnName, operator:"", value:""});
            this.selectStore.updateWhere(
                token, db.name, db.selectedTable,
                wheres,
            )
            this.whereCol = ""
            this.whereOperator = ""
            this.whereVal = ""
        },

        addWhereOp(operator: string) {

            const wheres = this.wheres;

            wheres.push({columnName: "", operator, value:""});
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.selectStore.updateWhere(
                token, db.name, db.selectedTable,
                wheres,
            )
            this.whereCol = ""
            this.whereOperator = ""
            this.whereVal = ""
        },

        addWhereVal(value: string) {

            const wheres = this.wheres;

            wheres.push({columnName: "", operator:"", value});
            const token = this.authStore.mustGetCurrent;
            const db = this.schemaStore.mustGetCurrentDatabase(token);
            this.selectStore.updateWhere(
                token, db.name, db.selectedTable,
                wheres,
            )
            this.whereCol = ""
            this.whereOperator = ""
            this.whereVal = ""
        },

    }
}
</script>

<template>
        <div class="row g-1">
            <div :class="{
                'col-1': !fieldsetSelectOpen,
                'col-4': fieldsetSelectOpen,
            }">
                <fieldset>
                    <legend @click="openFieldsetSelect" style="cursor:pointer;">
                        Select
                        <template v-if="columns.length > 0">
                            ({{ columns.length }})
                        </template>
                    </legend>
                    <template v-if="fieldsetSelectOpen">
                        <template v-for="(column, i) in columns" :key="i">
                            <div class="input-group mb-1">
                                <select 
                                  class="form-control"
                                  @change="setColumnFunc(i, (($event as InputEvent).target as HTMLSelectElement).value)" 
                                  :value="column.func">
                                    <option></option>
                                    <optgroup label="Functions">
                                        <option>char_length</option>
                                        <option>date</option>
                                        <option>from_unixtime</option>
                                        <option>lower</option>
                                        <option>round</option>
                                        <option>floor</option>
                                        <option>ceil</option>
                                        <option>sec_to_time</option>
                                        <option>time_to_sec</option>
                                        <option>upper</option>
                                    </optgroup>
                                    <optgroup label="Aggregation">
                                        <option>distinct</option>
                                        <option>avg</option>
                                        <option>count</option>
                                        <option>count distinct</option>
                                        <option>group_concat</option>
                                        <option>max</option>
                                        <option>min</option>
                                        <option>sum</option>
                                    </optgroup>
                                </select>
                                <span class="input-group-text">(</span>
                                <select class="form-control" @change="setColumnName(i, (($event as InputEvent).target as HTMLSelectElement).value)" :value="column.name">
                                    <option value=""></option>
                                    <option :value="col.name" v-for="col in table!.columns">{{ col.name }}</option>
                                </select>
                                <span class="input-group-text">)</span>
                                <button class="btn btn-danger" @click="removeColumn(i)">-</button>
                              </div>
                            </div>
                        </template>
                        <div class="input-group mb-1">
                            <select class="form-control" @change="addColumnFunc((($event as InputEvent).target as HTMLSelectElement).value)" v-model="selectFunc">
                                <option></option>
                                <optgroup label="Functions">
                                    <option>char_length</option>
                                    <option>date</option>
                                    <option>from_unixtime</option>
                                    <option>lower</option>
                                    <option>round</option>
                                    <option>floor</option>
                                    <option>ceil</option>
                                    <option>sec_to_time</option>
                                    <option>time_to_sec</option>
                                    <option>upper</option>
                                </optgroup>
                                <optgroup label="Aggregation">
                                    <option>distinct</option>
                                    <option>avg</option>
                                    <option>count</option>
                                    <option>count distinct</option>
                                    <option>group_concat</option>
                                    <option>max</option>
                                    <option>min</option>
                                    <option>sum</option>
                                </optgroup>
                            </select>
                            <span class="input-group-text">(</span>
                            <select class="form-control" @change="addColumnName((($event as InputEvent).target as HTMLSelectElement).value)" v-model="selectName">
                                <option value=""></option>
                                <option :value="col.name" v-for="col in table!.columns">{{ col.name }}</option>
                            </select>
                            <span class="input-group-text">)</span>
                            <button class="btn btn-danger disabled">-</button>
                        </div>
                    </template>
                </fieldset>
            </div>


            <div :class="{
                'col-1': !fieldsetSearchOpen,
                'col-4': fieldsetSearchOpen,
            }">
                <fieldset>
                    <legend @click="openFieldsetSearch" style="cursor:pointer;">
                        Search
                        <template v-if="wheres.length > 0">
                            ({{ wheres.length }})
                        </template>
                    </legend>
                    <template v-if="fieldsetSearchOpen">
                        <template v-for="(where, i) in wheres">
                        <div class="input-group mb-1">
                            <select class="form-control" v-model="where.columnName">
                                <option value=""></option>
                                <option :value="col.name" v-for="col in table!.columns">{{ col.name }}</option>
                            </select>
                            <select class="form-control" v-model="where.operator">
                                <option>=</option>
                                <option>&lt;</option>
                                <option>&gt;</option>
                                <option>&lt;=</option>
                                <option>&gt;=</option>
                                <option>!=</option>
                                <option>LIKE</option>
                                <option>LIKE %%</option>
                                <option>SQL</option>
                                <option>REGEXP</option>
                                <option>IN</option>
                                <option>FIND_IN_SET</option>
                                <option>IS NULL</option>
                                <option>NOT LIKE</option>
                                <option>NOT REGEXP</option>
                                <option>NOT IN</option>
                                <option>IS NOT NULL</option>
                            </select>
                          <input type="text" class="form-control" v-model="where.value"/>
                          <button class="btn btn-danger"  @click="removeWhere(i)">-</button>
                        </div>
                        </template>
                        <div class="input-group mb-1">
                            <select class="form-control" @change="addWhereCol((($event as InputEvent).target as HTMLSelectElement).value)" v-model="whereCol">
                                <option value=""></option>
                                <option :value="col.name" v-for="col in table!.columns">{{ col.name }}</option>
                            </select>
                            <select class="form-control" @change="addWhereOp((($event as InputEvent).target as HTMLSelectElement).value)" v-model="whereOperator">
                                <option>=</option>
                                <option>&lt;</option>
                                <option>&gt;</option>
                                <option>&lt;=</option>
                                <option>&gt;=</option>
                                <option>!=</option>
                                <option>LIKE</option>
                                <option>LIKE %%</option>
                                <option>SQL</option>
                                <option>REGEXP</option>
                                <option>IN</option>
                                <option>FIND_IN_SET</option>
                                <option>IS NULL</option>
                                <option>NOT LIKE</option>
                                <option>NOT REGEXP</option>
                                <option>NOT IN</option>
                                <option>IS NOT NULL</option>
                            </select>
                            <input type="text" class="form-control" @change="addWhereVal((($event as InputEvent).target as HTMLSelectElement).value)" v-model="whereVal"/>
                            <button class="btn btn-danger disabled">-</button>
                        </div>
                    </template>
                </fieldset>
            </div>


            <div :class="{
                'col-1': !fieldsetScopeOpen,
                'col-3': fieldsetScopeOpen,
            }">
                <fieldset>
                    <legend @click="openFieldsetScope" style="cursor:pointer;">
                        Scope
                    </legend>
                    <template v-if="fieldsetScopeOpen">

                        <div class="form-group">
                            <div class="input-group">
                                <span class="input-group-text">
                                    <label for="scope-limit">Limit</label>
                                </span>
                                <input type="number" name="scope-limit" class="form-control" :value="limit" @input="setLimit((($event as InputEvent).target as HTMLSelectElement).value)"/>
                            </div>
                        </div>
                        <div class="form-group">
                            <div class="input-group">
                                <span class="input-group-text">
                                    <label for="scope-offset">Offset</label>
                                </span>
                                <input type="number" name="scope-offset" class="form-control" :value="offset" @input="setOffset((($event as InputEvent).target as HTMLSelectElement).value)"/>
                            </div>
                        </div>
                    </template>
                </fieldset>
            </div>

            <div :class="{
                'col-1': !fieldsetSortOpen,
                'col-3': fieldsetSortOpen,
            }">
                <fieldset>
                    <legend @click="openFieldsetSort" style="cursor:pointer;">
                        Sort
                    </legend>
                    <template v-if="fieldsetSortOpen">
                        <template v-for="(where, i) in wheres">
                            <div class="form-group">2
                                <div class="input-group">
                                    <span class="input-group-text">
                                        <label for="scope-limit">Limit</label>
                                    </span>
                                    <select class="form-control" v-model="where.columnName">
                                    <option value=""></option>
                                    <option :value="col.name" v-for="col in table!.columns">{{ col.name }}</option>
                                </select>
                                </div>
                            </div>
                        <div class="form-group">
                            <div class="input-group">
                                <span class="input-group-text">
                                    <label for="scope-offset">Offset</label>
                                </span>
                                <input type="number" name="scope-offset" class="form-control" :value="offset" @input="setOffset((($event as InputEvent).target as HTMLSelectElement).value)"/>
                            </div>
                        </div>
                    </template>
                </fieldset>
            </div>

            <div class="col-1">
                <fieldset>
                    <legend>
                        Action
                    </legend>
                    <button class="btn btn-success" @click="select">
                        Select
                    </button>
                </fieldset>
            </div>

        </div>
</template>
