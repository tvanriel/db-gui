<script lang="ts">
import { PropType } from 'vue';
import { useAuthStore } from '../store/store';

export type BreadcrumbItem = {
    text: string,
    href?: string,
    active?: boolean,
}

export default {
    props: {
        items: {
            type: Array,
            validator: (list: BreadcrumbItem[]) => {
                return Array.prototype.every.call(list, (item: BreadcrumbItem) => {
                    return typeof item.active === "boolean" &&
                    (
                        typeof item.href === "undefined" ||
                        typeof item.href === "string"
                    ) && 
                    typeof item.text === "string"
                })
            },
            required: true,
        } as unknown as PropType<BreadcrumbItem[]>,
    },
    setup() {
        return {
            authStore: useAuthStore(),
        }
    },

    computed: {
        borderColour() {
            if (!this.authStore.hasCurrentConnection) {
                return '';
            }
            return 'border-' + this.authStore.currentConnection?.colour;
        }
    }
}
</script>

<template>
    <ol :class="'breadcrumb border ' + borderColour">
        <li v-for="item in items"
            :class="{'breadcrumb-item': true, 'active': item.active}">
            <template v-if="(typeof item.href === 'string' && item.active !== true)">
                <router-link :to="(item.href as string)">
                    {{item.text}}
                </router-link>
            </template>
            <span v-else="">
                {{item.text}}
            </span>
        </li>
    </ol>
</template>