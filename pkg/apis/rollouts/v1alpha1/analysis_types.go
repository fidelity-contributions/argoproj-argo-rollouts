package v1alpha1

import (
	"encoding/json"
	"time"

	intstrutil "k8s.io/apimachinery/pkg/util/intstr"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterAnalysisTemplate holds the template for performing canary analysis
// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=clusteranalysistemplates,shortName=cat,scope=Cluster
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time since resource was created"
type ClusterAnalysisTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec AnalysisTemplateSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// AnalysisTemplateList is a list of AnalysisTemplate resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterAnalysisTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ClusterAnalysisTemplate `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// AnalysisTemplate holds the template for performing canary analysis
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=analysistemplates,shortName=at
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time since resource was created"
type AnalysisTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec AnalysisTemplateSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// AnalysisTemplateList is a list of AnalysisTemplate resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AnalysisTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []AnalysisTemplate `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// AnalysisTemplateSpec is the specification for a AnalysisTemplate resource
type AnalysisTemplateSpec struct {
	// Metrics contains the list of metrics to query as part of an analysis run
	// +patchMergeKey=name
	// +patchStrategy=merge
	Metrics []Metric `json:"metrics,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,1,rep,name=metrics"`
	// Args are the list of arguments to the template
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +optional
	Args []Argument `json:"args,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=args"`
	// DryRun object contains the settings for running the analysis in Dry-Run mode
	// +patchMergeKey=metricName
	// +patchStrategy=merge
	// +optional
	DryRun []DryRun `json:"dryRun,omitempty" patchStrategy:"merge" patchMergeKey:"metricName" protobuf:"bytes,3,rep,name=dryRun"`
	// MeasurementRetention object contains the settings for retaining the number of measurements during the analysis
	// +patchMergeKey=metricName
	// +patchStrategy=merge
	// +optional
	MeasurementRetention []MeasurementRetention `json:"measurementRetention,omitempty" patchStrategy:"merge" patchMergeKey:"metricName" protobuf:"bytes,4,rep,name=measurementRetention"`
	// Templates reference to a list of analysis templates to combine with the rest of the metrics for an AnalysisRun
	// +patchMergeKey=templateName
	// +patchStrategy=merge
	Templates []AnalysisTemplateRef `json:"templates,omitempty" patchStrategy:"merge" patchMergeKey:"templateName" protobuf:"bytes,5,rep,name=templates"`
}

// DurationString is a string representing a duration (e.g. 30s, 5m, 1h)
type DurationString string

// Duration converts DurationString into a time.Duration
func (d DurationString) Duration() (time.Duration, error) {
	return time.ParseDuration(string(d))
}

// Metric defines a metric in which to perform analysis
type Metric struct {
	// Name is the name of the metric
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Interval defines an interval string (e.g. 30s, 5m, 1h) between each measurement.
	// If omitted, will perform a single measurement
	Interval DurationString `json:"interval,omitempty" protobuf:"bytes,2,opt,name=interval,casttype=DurationString"`
	// InitialDelay how long the AnalysisRun should wait before starting this metric
	InitialDelay DurationString `json:"initialDelay,omitempty" protobuf:"bytes,3,opt,name=initialDelay,casttype=DurationString"`
	// Count is the number of times to run the measurement. If both interval and count are omitted,
	// the effective count is 1. If only interval is specified, metric runs indefinitely.
	// If count > 1, interval must be specified.
	Count *intstrutil.IntOrString `json:"count,omitempty" protobuf:"bytes,4,opt,name=count"`
	// SuccessCondition is an expression which determines if a measurement is considered successful
	// Expression is a goevaluate expression. The keyword `result` is a variable reference to the
	// value of measurement. Results can be both structured data or primitive.
	// Examples:
	//   result > 10
	//   (result.requests_made * result.requests_succeeded / 100) >= 90
	SuccessCondition string `json:"successCondition,omitempty" protobuf:"bytes,5,opt,name=successCondition"`
	// FailureCondition is an expression which determines if a measurement is considered failed
	// If both success and failure conditions are specified, and the measurement does not fall into
	// either condition, the measurement is considered Inconclusive
	FailureCondition string `json:"failureCondition,omitempty" protobuf:"bytes,6,opt,name=failureCondition"`
	// FailureLimit is the maximum number of times the measurement is allowed to fail, before the
	// entire metric is considered Failed (default: 0)
	// -1 for making it disabled (when opting to use ConsecutiveSuccessLimit solely)
	FailureLimit *intstrutil.IntOrString `json:"failureLimit,omitempty" protobuf:"bytes,7,opt,name=failureLimit"`
	// InconclusiveLimit is the maximum number of times the measurement is allowed to measure
	// Inconclusive, before the entire metric is considered Inconclusive (default: 0)
	InconclusiveLimit *intstrutil.IntOrString `json:"inconclusiveLimit,omitempty" protobuf:"bytes,8,opt,name=inconclusiveLimit"`
	// ConsecutiveErrorLimit is the maximum number of times the measurement is allowed to error in
	// succession, before the metric is considered error (default: 4)
	ConsecutiveErrorLimit *intstrutil.IntOrString `json:"consecutiveErrorLimit,omitempty" protobuf:"bytes,9,opt,name=consecutiveErrorLimit"`
	// Provider configuration to the external system to use to verify the analysis
	Provider MetricProvider `json:"provider" protobuf:"bytes,10,opt,name=provider"`
	// ConsecutiveSuccessLimit is the number of consecutive times the measurement must succeed for the
	// entire metric to be considered Successful (default: 0, which means it's disabled)
	ConsecutiveSuccessLimit *intstrutil.IntOrString `json:"consecutiveSuccessLimit,omitempty" protobuf:"bytes,11,opt,name=consecutiveSuccessLimit"`
}

// DryRun defines the settings for running the analysis in Dry-Run mode.
type DryRun struct {
	// Name of the metric which needs to be evaluated in the Dry-Run mode. Wildcard '*' is supported and denotes all
	// the available metrics.
	MetricName string `json:"metricName" protobuf:"bytes,1,opt,name=metricName"`
}

// MeasurementRetention defines the settings for retaining the number of measurements during the analysis.
type MeasurementRetention struct {
	// MetricName is the name of the metric on which this retention policy should be applied.
	MetricName string `json:"metricName" protobuf:"bytes,1,opt,name=metricName"`
	// Limit is the maximum number of measurements to be retained for this given metric.
	Limit int32 `json:"limit" protobuf:"varint,2,opt,name=limit"`
}

// TTLStrategy defines the strategy for the time to live depending on if the analysis succeeded or failed
type TTLStrategy struct {
	// SecondsAfterCompletion is the number of seconds to live after completion.
	SecondsAfterCompletion *int32 `json:"secondsAfterCompletion,omitempty" protobuf:"varint,1,opt,name=secondsAfterCompletion"`
	// SecondsAfterFailure is the number of seconds to live after failure.
	SecondsAfterFailure *int32 `json:"secondsAfterFailure,omitempty" protobuf:"varint,2,opt,name=secondsAfterFailure"`
	// SecondsAfterSuccess is the number of seconds to live after success.
	SecondsAfterSuccess *int32 `json:"secondsAfterSuccess,omitempty" protobuf:"varint,3,opt,name=secondsAfterSuccess"`
}

// EffectiveCount is the effective count based on whether or not count/interval is specified
// If neither count or interval is specified, the effective count is 1
// If only interval is specified, metric runs indefinitely and there is no effective count (nil)
// Otherwise, it is the user specified value
func (m *Metric) EffectiveCount() *intstrutil.IntOrString {
	// Need to check if type is String
	if m.Count == nil || m.Count.IntValue() == 0 {
		if m.Interval == "" {
			one := intstrutil.FromInt(1)
			return &one
		}
		return nil
	}
	return m.Count
}

// MetricProvider which external system to use to verify the analysis
// Only one of the fields in this struct should be non-nil
type MetricProvider struct {
	// Prometheus specifies the prometheus metric to query
	Prometheus *PrometheusMetric `json:"prometheus,omitempty" protobuf:"bytes,1,opt,name=prometheus"`
	// Kayenta specifies a Kayenta metric
	Kayenta *KayentaMetric `json:"kayenta,omitempty" protobuf:"bytes,2,opt,name=kayenta"`
	// Web specifies a generic HTTP web metric
	Web *WebMetric `json:"web,omitempty" protobuf:"bytes,3,opt,name=web"`
	// Datadog specifies a datadog metric to query
	Datadog *DatadogMetric `json:"datadog,omitempty" protobuf:"bytes,4,opt,name=datadog"`
	// Wavefront specifies the wavefront metric to query
	Wavefront *WavefrontMetric `json:"wavefront,omitempty" protobuf:"bytes,5,opt,name=wavefront"`
	// NewRelic specifies the newrelic metric to query
	NewRelic *NewRelicMetric `json:"newRelic,omitempty" protobuf:"bytes,6,opt,name=newRelic"`
	// Job specifies the job metric run
	Job *JobMetric `json:"job,omitempty" protobuf:"bytes,7,opt,name=job"`
	// CloudWatch specifies the cloudWatch metric to query
	CloudWatch *CloudWatchMetric `json:"cloudWatch,omitempty" protobuf:"bytes,8,opt,name=cloudWatch"`
	// Graphite specifies the Graphite metric to query
	Graphite *GraphiteMetric `json:"graphite,omitempty" protobuf:"bytes,9,opt,name=graphite"`
	// Influxdb specifies the influxdb metric to query
	Influxdb *InfluxdbMetric `json:"influxdb,omitempty" protobuf:"bytes,10,opt,name=influxdb"`
	// SkyWalking specifies the skywalking metric to query
	SkyWalking *SkyWalkingMetric `json:"skywalking,omitempty" protobuf:"bytes,11,opt,name=skywalking"`
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Type=object
	// Plugin specifies the hashicorp go-plugin metric to query
	Plugin map[string]json.RawMessage `json:"plugin,omitempty" protobuf:"bytes,12,opt,name=plugin"`
}

// AnalysisPhase is the overall phase of an AnalysisRun, MetricResult, or Measurement
type AnalysisPhase string

// Possible AnalysisPhase values
const (
	AnalysisPhasePending      AnalysisPhase = "Pending"
	AnalysisPhaseRunning      AnalysisPhase = "Running"
	AnalysisPhaseSuccessful   AnalysisPhase = "Successful"
	AnalysisPhaseFailed       AnalysisPhase = "Failed"
	AnalysisPhaseError        AnalysisPhase = "Error"
	AnalysisPhaseInconclusive AnalysisPhase = "Inconclusive"
)

// Completed returns whether or not the analysis status is considered completed
func (as AnalysisPhase) Completed() bool {
	switch as {
	case AnalysisPhaseSuccessful, AnalysisPhaseFailed, AnalysisPhaseError, AnalysisPhaseInconclusive:
		return true
	}
	return false
}

// Arguments to perform a prometheus range query
type PrometheusRangeQueryArgs struct {
	// The start time to query in expr format e.g. now(), now() - duration("1h"), now() - duration("{{args.lookback_duration}}")
	Start string `json:"start,omitempty" protobuf:"bytes,1,opt,name=start"`
	// The end time to query in expr format e.g. now(), now() - duration("1h"), now() - duration("{{args.lookback_duration}}")
	End string `json:"end,omitempty" protobuf:"bytes,2,opt,name=end"`
	// The maximum time between two slices from the start to end (e.g. 30s, 5m, 1h).
	Step DurationString `json:"step,omitempty" protobuf:"bytes,3,opt,name=step,casttype=DurationString"`
}

// PrometheusMetric defines the prometheus query to perform canary analysis
type PrometheusMetric struct {
	// Address is the HTTP address and port of the prometheus server
	Address string `json:"address,omitempty" protobuf:"bytes,1,opt,name=address"`
	// Query is a raw prometheus query to perform
	Query string `json:"query,omitempty" protobuf:"bytes,2,opt,name=query"`
	// Authentication details
	// +optional
	Authentication Authentication `json:"authentication,omitempty" protobuf:"bytes,3,opt,name=authentication"`
	// Timeout represents the duration within which a prometheus query should complete. It is expressed in seconds.
	// +optional
	Timeout *int64 `json:"timeout,omitempty" protobuf:"bytes,4,opt,name=timeout"`
	// Insecure skips host TLS verification
	Insecure bool `json:"insecure,omitempty" protobuf:"varint,5,opt,name=insecure"`
	// Headers are optional HTTP headers to use in the request
	// +optional
	// +patchMergeKey=key
	// +patchStrategy=merge
	Headers []WebMetricHeader `json:"headers,omitempty" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,6,opt,name=headers"`
	// Arguments for prometheus
	// +optional
	RangeQuery *PrometheusRangeQueryArgs `json:"rangeQuery,omitempty" protobuf:"bytes,7,opt,name=rangeQuery"`
}

// Authentication method
type Authentication struct {
	// Sigv4 Config is the aws SigV4 configuration to use for SigV4 signing if using Amazon Managed Prometheus
	// +optional
	Sigv4 Sigv4Config `json:"sigv4,omitempty" protobuf:"bytes,1,opt,name=sigv4"`
	// OAuth2 config
	// +optional
	OAuth2 OAuth2Config `json:"oauth2,omitempty" protobuf:"bytes,2,opt,name=oauth2"`
}

type OAuth2Config struct {
	// OAuth2 provider token URL
	TokenURL string `json:"tokenUrl,omitempty" protobuf:"bytes,1,name=tokenUrl"`
	// OAuth2 client ID
	ClientID string `json:"clientId,omitempty" protobuf:"bytes,2,name=clientId"`
	// OAuth2 client secret
	ClientSecret string `json:"clientSecret,omitempty" protobuf:"bytes,3,name=clientSecret"`
	// OAuth2 scopes
	// +optional
	Scopes []string `json:"scopes,omitempty" protobuf:"bytes,4,opt,name=scopes"`
}

type Sigv4Config struct {
	// Region is the AWS Region to sign the SigV4 Request
	Region string `json:"region,omitempty" protobuf:"bytes,1,opt,name=address"`
	// Profile is the Credential Profile used to sign the SigV4 Request
	Profile string `json:"profile,omitempty" protobuf:"bytes,2,opt,name=profile"`
	// RoleARN is the IAM role used to sign the SIgV4 Request
	RoleARN string `json:"roleArn,omitempty" protobuf:"bytes,3,opt,name=roleArn"`
}

// WavefrontMetric defines the wavefront query to perform canary analysis
type WavefrontMetric struct {
	// Address is the HTTP address and port of the wavefront server
	Address string `json:"address,omitempty" protobuf:"bytes,1,opt,name=address"`
	// Query is a raw wavefront query to perform
	Query string `json:"query,omitempty" protobuf:"bytes,2,opt,name=query"`
}

// NewRelicMetric defines the newrelic query to perform canary analysis
type NewRelicMetric struct {
	// Profile is the name of the secret holding NR account configuration
	Profile string `json:"profile,omitempty" protobuf:"bytes,1,opt,name=profile"`
	// Query is a raw newrelic NRQL query to perform
	Query string `json:"query" protobuf:"bytes,2,opt,name=query"`
	// Timeout represents the duration limit in seconds that will apply to the NRQL query
	// +optional
	Timeout *int64 `json:"timeout,omitempty" protobuf:"bytes,3,opt,name=timeout"`
}

// JobMetric defines a job to run which acts as a metric
type JobMetric struct {
	Metadata metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec     batchv1.JobSpec   `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// GraphiteMetric defines the Graphite query to perform canary analysis
type GraphiteMetric struct {
	// Address is the HTTP address and port of the Graphite server
	Address string `json:"address,omitempty" protobuf:"bytes,1,opt,name=address"`
	// Query is a raw Graphite query to perform
	Query string `json:"query,omitempty" protobuf:"bytes,2,opt,name=query"`
}

// InfluxdbMetric defines the InfluxDB Flux query to perform canary analysis
type InfluxdbMetric struct {
	// Profile is the name of the secret holding InfluxDB account configuration
	Profile string `json:"profile,omitempty" protobuf:"bytes,1,opt,name=profile"`
	// Query is a raw InfluxDB flux query to perform
	Query string `json:"query,omitempty" protobuf:"bytes,2,opt,name=query"`
}

// CloudWatchMetric defines the cloudwatch query to perform canary analysis
type CloudWatchMetric struct {
	Interval          DurationString              `json:"interval,omitempty" protobuf:"bytes,1,opt,name=interval,casttype=DurationString"`
	MetricDataQueries []CloudWatchMetricDataQuery `json:"metricDataQueries" protobuf:"bytes,2,rep,name=metricDataQueries"`
}

// CloudWatchMetricDataQuery defines the cloudwatch query
type CloudWatchMetricDataQuery struct {
	Id         string                  `json:"id,omitempty" protobuf:"bytes,1,opt,name=id"`
	Expression *string                 `json:"expression,omitempty" protobuf:"bytes,2,opt,name=expression"`
	Label      *string                 `json:"label,omitempty" protobuf:"bytes,3,opt,name=label"`
	MetricStat *CloudWatchMetricStat   `json:"metricStat,omitempty" protobuf:"bytes,4,opt,name=metricStat"`
	Period     *intstrutil.IntOrString `json:"period,omitempty" protobuf:"varint,5,opt,name=period"`
	ReturnData *bool                   `json:"returnData,omitempty" protobuf:"bytes,6,opt,name=returnData"`
}

type CloudWatchMetricStat struct {
	Metric CloudWatchMetricStatMetric `json:"metric,omitempty" protobuf:"bytes,1,opt,name=metric"`
	Period intstrutil.IntOrString     `json:"period,omitempty" protobuf:"varint,2,opt,name=period"`
	Stat   string                     `json:"stat,omitempty" protobuf:"bytes,3,opt,name=stat"`
	Unit   string                     `json:"unit,omitempty" protobuf:"bytes,4,opt,name=unit"`
}

type CloudWatchMetricStatMetric struct {
	Dimensions []CloudWatchMetricStatMetricDimension `json:"dimensions,omitempty" protobuf:"bytes,1,rep,name=dimensions"`
	MetricName string                                `json:"metricName,omitempty" protobuf:"bytes,2,opt,name=metricName"`
	Namespace  *string                               `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}

type CloudWatchMetricStatMetricDimension struct {
	Name  string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
}

// AnalysisRun is an instantiation of an AnalysisTemplate
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=analysisruns, shortName=ar
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase",description="AnalysisRun status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time since resource was created"
type AnalysisRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              AnalysisRunSpec   `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status            AnalysisRunStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// AnalysisRunList is a list of AnalysisTemplate resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AnalysisRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []AnalysisRun `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type SkyWalkingMetric struct {
	Address  string         `json:"address,omitempty" protobuf:"bytes,1,opt,name=address"`
	Query    string         `json:"query,omitempty" protobuf:"bytes,2,opt,name=query"`
	Interval DurationString `json:"interval,omitempty" protobuf:"bytes,3,opt,name=interval,casttype=DurationString"`
}

// AnalysisRunSpec is the spec for a AnalysisRun resource
type AnalysisRunSpec struct {
	// Metrics contains the list of metrics to query as part of an analysis run
	// +patchMergeKey=name
	// +patchStrategy=merge
	Metrics []Metric `json:"metrics" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,1,rep,name=metrics"`
	// Args are the list of arguments used in this run
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Args []Argument `json:"args,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=args"`
	// Terminate is used to prematurely stop the run (e.g. rollout completed and analysis is no longer desired)
	Terminate bool `json:"terminate,omitempty" protobuf:"varint,3,opt,name=terminate"`
	// DryRun object contains the settings for running the analysis in Dry-Run mode
	// +patchMergeKey=metricName
	// +patchStrategy=merge
	// +optional
	DryRun []DryRun `json:"dryRun,omitempty" patchStrategy:"merge" patchMergeKey:"metricName" protobuf:"bytes,4,rep,name=dryRun"`
	// MeasurementRetention object contains the settings for retaining the number of measurements during the analysis
	// +patchMergeKey=metricName
	// +patchStrategy=merge
	// +optional
	MeasurementRetention []MeasurementRetention `json:"measurementRetention,omitempty" patchStrategy:"merge" patchMergeKey:"metricName" protobuf:"bytes,5,rep,name=measurementRetention"`
	// TTLStrategy object contains the strategy for the time to live depending on if the analysis succeeded or failed
	// +optional
	TTLStrategy *TTLStrategy `json:"ttlStrategy,omitempty" protobuf:"bytes,6,opt,name=ttlStrategy"`
}

// Argument is an argument to an AnalysisRun
type Argument struct {
	// Name is the name of the argument
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Value is the value of the argument
	// +optional
	Value *string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
	// ValueFrom is a reference to where a secret is stored. This field is one of the fields with valueFrom
	// +optional
	ValueFrom *ValueFrom `json:"valueFrom,omitempty" protobuf:"bytes,3,opt,name=valueFrom"`
}

type ValueFrom struct {
	// Secret is a reference to where a secret is stored. This field is one of the fields with valueFrom
	// +optional
	SecretKeyRef *SecretKeyRef `json:"secretKeyRef,omitempty" protobuf:"bytes,1,opt,name=secretKeyRef"`
	// FieldRef is a reference to the fields in metadata which we are referencing. This field is one of the fields with
	// valueFrom
	// +optional
	FieldRef *FieldRef `json:"fieldRef,omitempty" protobuf:"bytes,2,opt,name=fieldRef"`
}

type SecretKeyRef struct {
	// Name is the name of the secret
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Key is the key of the secret to select from.
	Key string `json:"key" protobuf:"bytes,2,opt,name=key"`
}

// AnalysisRunStatus is the status for a AnalysisRun resource
type AnalysisRunStatus struct {
	// Phase is the status of the analysis run
	Phase AnalysisPhase `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=AnalysisPhase"`
	// Message is a message explaining current status
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
	// MetricResults contains the metrics collected during the run
	MetricResults []MetricResult `json:"metricResults,omitempty" protobuf:"bytes,3,rep,name=metricResults"`
	// StartedAt indicates when the analysisRun first started
	StartedAt *metav1.Time `json:"startedAt,omitempty" protobuf:"bytes,4,opt,name=startedAt"`
	// RunSummary contains the final results from the metric executions
	RunSummary RunSummary `json:"runSummary,omitempty" protobuf:"bytes,5,opt,name=runSummary"`
	// DryRunSummary contains the final results from the metric executions in the dry-run mode
	DryRunSummary *RunSummary `json:"dryRunSummary,omitempty" protobuf:"bytes,6,opt,name=dryRunSummary"`
	// CompletedAt indicates when the analysisRun completed
	CompletedAt *metav1.Time `json:"completedAt,omitempty" protobuf:"bytes,7,opt,name=completedAt"`
}

// RunSummary contains the final results from the metric executions
type RunSummary struct {
	// This is equal to the sum of Successful, Failed, Inconclusive
	Count int32 `json:"count,omitempty" protobuf:"varint,1,opt,name=count"`
	// Successful is the number of times the metric was measured Successful
	Successful int32 `json:"successful,omitempty" protobuf:"varint,2,opt,name=successful"`
	// Failed is the number of times the metric was measured Failed
	Failed int32 `json:"failed,omitempty" protobuf:"varint,3,opt,name=failed"`
	// Inconclusive is the number of times the metric was measured Inconclusive
	Inconclusive int32 `json:"inconclusive,omitempty" protobuf:"varint,4,opt,name=inconclusive"`
	// Error is the number of times an error was encountered during measurement
	Error int32 `json:"error,omitempty" protobuf:"varint,5,opt,name=error"`
}

// MetricResult contain a list of the most recent measurements for a single metric along with
// counters on how often the measurement
type MetricResult struct {
	// Name is the name of the metric
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Phase is the overall aggregate status of the metric
	Phase AnalysisPhase `json:"phase" protobuf:"bytes,2,opt,name=phase,casttype=AnalysisPhase"`
	// Measurements holds the most recent measurements collected for the metric
	Measurements []Measurement `json:"measurements,omitempty" protobuf:"bytes,3,rep,name=measurements"`
	// Message contains a message describing current condition (e.g. error messages)
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	// Count is the number of times the metric was measured without Error
	// This is equal to the sum of Successful, Failed, Inconclusive
	Count int32 `json:"count,omitempty" protobuf:"varint,5,opt,name=count"`
	// Successful is the number of times the metric was measured Successful
	Successful int32 `json:"successful,omitempty" protobuf:"varint,6,opt,name=successful"`
	// Failed is the number of times the metric was measured Failed
	Failed int32 `json:"failed,omitempty" protobuf:"varint,7,opt,name=failed"`
	// Inconclusive is the number of times the metric was measured Inconclusive
	Inconclusive int32 `json:"inconclusive,omitempty" protobuf:"varint,8,opt,name=inconclusive"`
	// Error is the number of times an error was encountered during measurement
	Error int32 `json:"error,omitempty" protobuf:"varint,9,opt,name=error"`
	// ConsecutiveError is the number of times an error was encountered during measurement in succession
	// Resets to zero when non-errors are encountered
	ConsecutiveError int32 `json:"consecutiveError,omitempty" protobuf:"varint,10,opt,name=consecutiveError"`
	// DryRun indicates whether this metric is running in a dry-run mode or not
	DryRun bool `json:"dryRun,omitempty" protobuf:"varint,11,opt,name=dryRun"`
	// Metadata stores additional metadata about this metric. It is used by different providers to store
	// the final state which gets used while taking measurements. For example, Prometheus uses this field
	// to store the final resolved query after substituting the template arguments.
	Metadata map[string]string `json:"metadata,omitempty" protobuf:"bytes,12,rep,name=metadata"`
	// ConsecutiveSuccess is the number of times a measurement was successful in succession
	// Resets to zero when failures, inconclusive measurements, or errors are encountered
	ConsecutiveSuccess int32 `json:"consecutiveSuccess,omitempty" protobuf:"varint,13,opt,name=consecutiveSuccess"`
}

// Measurement is a point in time result value of a single metric, and the time it was measured
type Measurement struct {
	// Phase is the status of this single measurement
	Phase AnalysisPhase `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=AnalysisPhase"`
	// Message contains a message describing current condition (e.g. error messages)
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
	// StartedAt is the timestamp in which this measurement started to be measured
	StartedAt *metav1.Time `json:"startedAt,omitempty" protobuf:"bytes,3,opt,name=startedAt"`
	// FinishedAt is the timestamp in which this measurement completed and value was collected
	FinishedAt *metav1.Time `json:"finishedAt,omitempty" protobuf:"bytes,4,opt,name=finishedAt"`
	// Value is the measured value of the metric
	Value string `json:"value,omitempty" protobuf:"bytes,5,opt,name=value"`
	// Metadata stores additional metadata about this metric result, used by the different providers
	// (e.g. kayenta run ID, job name)
	Metadata map[string]string `json:"metadata,omitempty" protobuf:"bytes,6,rep,name=metadata"`
	// ResumeAt is the  timestamp when the analysisRun should try to resume the measurement
	ResumeAt *metav1.Time `json:"resumeAt,omitempty" protobuf:"bytes,7,opt,name=resumeAt"`
}

type KayentaMetric struct {
	Address string `json:"address" protobuf:"bytes,1,opt,name=address"`

	Application string `json:"application" protobuf:"bytes,2,opt,name=application"`

	CanaryConfigName string `json:"canaryConfigName" protobuf:"bytes,3,opt,name=canaryConfigName"`

	MetricsAccountName       string `json:"metricsAccountName" protobuf:"bytes,4,opt,name=metricsAccountName"`
	ConfigurationAccountName string `json:"configurationAccountName" protobuf:"bytes,5,opt,name=configurationAccountName"`
	StorageAccountName       string `json:"storageAccountName" protobuf:"bytes,6,opt,name=storageAccountName"`

	Threshold KayentaThreshold `json:"threshold" protobuf:"bytes,7,opt,name=threshold"`

	Scopes []KayentaScope `json:"scopes" protobuf:"bytes,8,rep,name=scopes"`
}

type KayentaThreshold struct {
	Pass     int64 `json:"pass" protobuf:"varint,1,opt,name=pass"`
	Marginal int64 `json:"marginal" protobuf:"varint,2,opt,name=marginal"`
}

type KayentaScope struct {
	Name            string      `json:"name" protobuf:"bytes,1,opt,name=name"`
	ControlScope    ScopeDetail `json:"controlScope" protobuf:"bytes,2,opt,name=controlScope"`
	ExperimentScope ScopeDetail `json:"experimentScope" protobuf:"bytes,3,opt,name=experimentScope"`
}

type ScopeDetail struct {
	Scope  string `json:"scope" protobuf:"bytes,1,opt,name=scope"`
	Region string `json:"region" protobuf:"bytes,2,opt,name=region"`
	Step   int64  `json:"step" protobuf:"varint,3,opt,name=step"`
	Start  string `json:"start" protobuf:"bytes,4,opt,name=start"`
	End    string `json:"end" protobuf:"bytes,5,opt,name=end"`
}

type WebMetric struct {
	// Method is the method of the web metric (empty defaults to GET)
	Method WebMetricMethod `json:"method,omitempty" protobuf:"bytes,1,opt,name=method"`
	// URL is the address of the web metric
	URL string `json:"url" protobuf:"bytes,2,opt,name=url"`
	// +patchMergeKey=key
	// +patchStrategy=merge
	// Headers are optional HTTP headers to use in the request
	Headers []WebMetricHeader `json:"headers,omitempty" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,3,rep,name=headers"`
	// Body is the body of the web metric (must be POST/PUT)
	Body string `json:"body,omitempty" protobuf:"bytes,4,opt,name=body"`
	// TimeoutSeconds is the timeout for the request in seconds (default: 10)
	TimeoutSeconds int64 `json:"timeoutSeconds,omitempty" protobuf:"varint,5,opt,name=timeoutSeconds"`
	// JSONPath is a JSON Path to use as the result variable (default: "{$}")
	JSONPath string `json:"jsonPath,omitempty" protobuf:"bytes,6,opt,name=jsonPath"`
	// Insecure skips host TLS verification
	Insecure bool `json:"insecure,omitempty" protobuf:"varint,7,opt,name=insecure"`
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Type=object
	// JSONBody is the body of the web metric in a json format (method must be POST/PUT)
	JSONBody json.RawMessage `json:"jsonBody,omitempty" protobuf:"bytes,8,opt,name=jsonBody,casttype=encoding/json.RawMessage"`
	// Authentication details
	// +optional
	Authentication Authentication `json:"authentication,omitempty" protobuf:"bytes,9,opt,name=authentication"`
}

// WebMetricMethod is the available HTTP methods
type WebMetricMethod string

// Possible HTTP method values
const (
	WebMetricMethodGet  WebMetricMethod = "GET"
	WebMetricMethodPost WebMetricMethod = "POST"
	WebMetricMethodPut  WebMetricMethod = "PUT"
)

type WebMetricHeader struct {
	Key   string `json:"key" protobuf:"bytes,1,opt,name=key"`
	Value string `json:"value" protobuf:"bytes,2,opt,name=value"`
}

type DatadogMetric struct {
	// +kubebuilder:default="5m"
	// Interval refers to the Interval time window in Datadog (default: 5m). Not to be confused with the polling rate for the metric.
	Interval DurationString `json:"interval,omitempty" protobuf:"bytes,1,opt,name=interval,casttype=DurationString"`
	Query    string         `json:"query,omitempty" protobuf:"bytes,2,opt,name=query"`
	// Queries is a map of query_name_as_key: query. You can then use query_name_as_key inside Formula.Used for v2
	// +kubebuilder:validation:Type=object
	Queries map[string]string `json:"queries,omitempty" protobuf:"bytes,3,opt,name=queries"`
	// Formula refers to the Formula made up of the queries. Only useful with Queries. Used for v2
	Formula string `json:"formula,omitempty" protobuf:"bytes,4,opt,name=formula"`
	// ApiVersion refers to the Datadog API version being used (default: v1). v1 will eventually be deprecated.
	// +kubebuilder:validation:Enum=v1;v2
	// +kubebuilder:default=v1
	ApiVersion string `json:"apiVersion,omitempty" protobuf:"bytes,5,opt,name=apiVersion"`
	// +kubebuilder:validation:Enum=avg;min;max;sum;last;percentile;mean;l2norm;area
	// Aggregator is a type of aggregator to use for metrics-based queries (default: ""). Used for v2
	Aggregator string `json:"aggregator,omitempty" protobuf:"bytes,6,opt,name=aggregator"`
	// Secret refers to the name of the secret that should be used for an analysis and should exists in the namespace where the controller is.
	// +optional
	SecretRef SecretRef `json:"secretRef,omitempty" protobuf:"bytes,7,opt,name=secretRef"`
}

type SecretRef struct {
	// Name refers to the name of the secret that should be used to integrate with Datadog.
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// Namespaced indicates whether the secret is in the namespace where rollouts it installed or in the namespace where the metric was found
	Namespaced bool `json:"namespaced,omitempty" protobuf:"varint,2,opt,namespaced=dryRun"`
}
