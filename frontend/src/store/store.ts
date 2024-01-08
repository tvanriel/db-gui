import {backendApi} from '../lib/backend-api';

import {AuthService} from '../services/auth-service';
import {ExportService} from '../services/export-service';
import {ImportService} from '../services/import-service';
import {SchemaService} from '../services/schema-service';
import {SelectService} from '../services/select-service';
import {SqlService} from '../services/sql-service';

import {makeAuthStore} from './auth/auth-store';
import {makeDatabaseStore} from './database/database-store';
import {makeExportStore} from './export/export-store';
import {makeImportStore} from './import/import-store';
import {makeSelectStore} from './select/select-store';
import {makeSqlStore} from './sql/sql-store';

export const useAuthStore = makeAuthStore(new AuthService(backendApi));
export const useSchemaStore = makeDatabaseStore(new SchemaService(backendApi));
export const useImportStore = makeImportStore(new ImportService(backendApi));
export const useExportStore = makeExportStore(new ExportService(backendApi));
export const useSelectStore = makeSelectStore(new SelectService(backendApi));
export const useSqlStore = makeSqlStore(new SqlService(backendApi));
