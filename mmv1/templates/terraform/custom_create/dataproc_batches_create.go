userAgent, err := generateUserAgentString(d, config.userAgent)
if err != nil {
	return err
}

obj := make(map[string]interface{})

labelsProp, err := expandDataprocBatchesLabels(d.Get("labels"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
	obj["labels"] = labelsProp
}
runtimeConfigProp, err := expandDataprocBatchesRuntimeConfig(d.Get("runtime_config"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("runtime_config"); !isEmptyValue(reflect.ValueOf(runtimeConfigProp)) && (ok || !reflect.DeepEqual(v, runtimeConfigProp)) {
	obj["runtimeConfig"] = runtimeConfigProp
}
environmentConfigProp, err := expandDataprocBatchesEnvironmentConfig(d.Get("environment_config"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("environment_config"); !isEmptyValue(reflect.ValueOf(environmentConfigProp)) && (ok || !reflect.DeepEqual(v, environmentConfigProp)) {
	obj["environmentConfig"] = environmentConfigProp
}
pysparkBatchProp, err := expandDataprocBatchesPysparkBatch(d.Get("pyspark_batch"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("pyspark_batch"); !isEmptyValue(reflect.ValueOf(pysparkBatchProp)) && (ok || !reflect.DeepEqual(v, pysparkBatchProp)) {
	obj["pysparkBatch"] = pysparkBatchProp
}
sparkBatchProp, err := expandDataprocBatchesSparkBatch(d.Get("spark_batch"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("spark_batch"); !isEmptyValue(reflect.ValueOf(sparkBatchProp)) && (ok || !reflect.DeepEqual(v, sparkBatchProp)) {
	obj["sparkBatch"] = sparkBatchProp
}
sparkRBatchProp, err := expandDataprocBatchesSparkRBatch(d.Get("spark_r_batch"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("spark_r_batch"); !isEmptyValue(reflect.ValueOf(sparkRBatchProp)) && (ok || !reflect.DeepEqual(v, sparkRBatchProp)) {
	obj["sparkRBatch"] = sparkRBatchProp
}
sparkSqlBatchProp, err := expandDataprocBatchesSparkSqlBatch(d.Get("spark_sql_batch"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("spark_sql_batch"); !isEmptyValue(reflect.ValueOf(sparkSqlBatchProp)) && (ok || !reflect.DeepEqual(v, sparkSqlBatchProp)) {
	obj["sparkSqlBatch"] = sparkSqlBatchProp
}

url, err := replaceVars(d, config, "{{DataprocBasePath}}projects/{{project}}/locations/{{location}}/batches?batchId={{batch_id}}")
if err != nil {
	return err
}

log.Printf("[DEBUG] Creating new Batches: %#v", obj)
billingProject := ""

project, err := getProject(d, config)
if err != nil {
	return fmt.Errorf("Error fetching project for Batches: %s", err)
}
billingProject = project

// err == nil indicates that the billing_project value was found
if bp, err := getBillingProject(d, config); err == nil {
	billingProject = bp
}

res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
if err != nil {
	return fmt.Errorf("Error creating Batches: %s", err)
}

// setting batch_id from the operation response
if err := d.Set("batch_id", flattenDataprocBatchesBatchId(res, d, config)); err != nil {
	log.Printf("[ERROR] Error while setting batch_id: %s", err)
}

// waiting for the end-state of batch execution
if err := waitForDataprocBatchState(d, config, billingProject, userAgent, d.Timeout(schema.TimeoutCreate), []string{"PENDING", "RUNNING", "CANCELLING"}); err != nil {
	return fmt.Errorf("Error while waiting for end state of batch with batch ID %q: %s", d.Id(), err)
}

// Store the ID now
id, err := replaceVars(d, config, "projects/{{project}}/locations/{{location}}/batches/{{batch_id}}")
if err != nil {
	return fmt.Errorf("Error constructing id: %s", err)
}
d.SetId(id)

log.Printf("[DEBUG] Finished creating Batches %q: %#v", d.Id(), res)

return resourceDataprocBatchesRead(d, meta)