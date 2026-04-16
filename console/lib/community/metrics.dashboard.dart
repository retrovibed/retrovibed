import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/community/api.dart' as communityapi;
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/design.kit/inputs.dart' as inputs;

List<timex.Range> _defaultSegments() {
  final now = DateTime.now();
  return [
    timex.Range(DateTime(now.year, now.month - 1, now.day), now),
    timex.Range(DateTime(now.year, now.month - 3, now.day), now),
    timex.Range(DateTime(now.year - 1, now.month, now.day), now),
    timex.Range(DateTime(now.year - 3, now.month, now.day), now),
  ];
}

class MetricsDashboard extends StatefulWidget {
  final Community community;
  final Future<MetricsSyncProgress> Function(
    String id, {
    List<httpx.Option> options,
  })
  apicommunitysync;
  final Future<CommunityMetricsResponse> Function(
    String id, {
    required DateTime startDate,
    required DateTime endDate,
    List<httpx.Option> options,
  })
  apicommunitymetrics;
  final List<timex.Range> segments;

  MetricsDashboard({
    super.key,
    required this.community,
    this.apicommunitysync = communityapi.API.sync,
    this.apicommunitymetrics = communityapi.API.metrics,
    List<timex.Range>? segments,
  }) : segments = segments ?? _defaultSegments();

  @override
  State<MetricsDashboard> createState() => _MetricsDashboardState();
}

class _MetricsDashboardState extends State<MetricsDashboard> {
  CommunityMetricsResponse _metrics = CommunityMetricsResponse();
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  timex.Range _selected = _defaultSegments().first;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _clearCause() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  void initState() {
    super.initState();
    _selected = widget.segments.first;
    _syncAndLoad();
  }

  Future<void> _syncAndLoad() {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    final auth = [authn.AuthzCache.bearer(context)];

    return httpx
        .withRetry(() => widget.apicommunitysync(widget.community.id, options: auth))
        .then((_) => _loadMetrics())
        .catchError((cause) {
          setState(() {
            _cause = ds.Errors.httpauto(cause, onTap: _clearCause);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _clearCause);
          });
        })
        .whenComplete(() {
          setState(() {
            _loading = false;
          });
        });
  }

  Future<void> _loadMetrics() {
    final auth = [authn.AuthzCache.bearer(context)];

    return httpx
        .withRetry(
          () => widget.apicommunitymetrics(
            widget.community.id,
            startDate: _selected.begin,
            endDate: _selected.end,
            options: auth,
          ),
        )
        .then((response) {
          setState(() {
            _metrics = response;
          });
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Errors.httpauto(cause, onTap: _clearCause);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _clearCause);
          });
        })
        .whenComplete(() {
          setState(() {
            _loading = false;
          });
        });
  }

  void _changeRange(timex.Range range) {
    setState(() {
      _selected = range;
    });
    _loadMetrics();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    final emptyState = Center(
      child: Column(
        spacing: defaults.spacing,
        children: [
          Icon(
            Icons.analytics_outlined,
            size: 48,
            color: theme.colorScheme.outline,
          ),
          Text(
            'No metrics data yet',
            style: theme.textTheme.bodyMedium?.copyWith(
              color: theme.colorScheme.outline,
            ),
          ),
        ],
      ),
    );

    final metricsContent =
        _metrics.hasSummary()
            ? Column(
              children: [
                Row(
                  children: [
                    Expanded(
                      child: _MetricCard(
                        title: 'Subscribers',
                        value: _metrics.summary.subscribers.toString(),
                        icon: Icons.people,
                      ),
                    ),
                    Expanded(
                      child: _MetricCard(
                        title: 'Archivers',
                        value: _metrics.totalArchivers.toString(),
                        icon: Icons.archive,
                      ),
                    ),
                  ],
                ),
                Visibility(
                  visible: _metrics.items.isNotEmpty,
                  child: Column(
                    children: [
                      SizedBox(height: defaults.spacing * 2),
                      Text(
                        'Per-Content Breakdown',
                        style: theme.textTheme.titleSmall,
                      ),
                      SizedBox(height: defaults.spacing),
                      ..._metrics.items.map((item) => _MetricRow(metric: item)),
                    ],
                  ),
                ),
              ],
            )
            : emptyState;

    return ds.Loading(
      loading: _loading,
      cause: _cause,
      Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Wrap(
            alignment: WrapAlignment.spaceBetween,
            crossAxisAlignment: WrapCrossAlignment.center,
            spacing: defaults.spacing,
            runSpacing: defaults.spacing,
            children: [
              Text('Metrics', style: theme.textTheme.titleMedium),
              ds.LoadingIconButton(
                icon: Icon(Icons.refresh),
                onPressed: _syncAndLoad,
                tooltip: 'Sync metrics',
              ),
            ],
          ),
          inputs.TimeRange(
            segments: widget.segments,
            selected: _selected,
            onChanged: _changeRange,
          ),
          metricsContent,
        ],
      ),
    );
  }
}

class _MetricCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;

  const _MetricCard({
    required this.title,
    required this.value,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return Card(
      child: Padding(
        padding: defaults.padding,
        child: Column(
          children: [
            Icon(icon, color: theme.colorScheme.primary),
            SizedBox(height: defaults.spacing / 2),
            Text(
              value,
              style: theme.textTheme.headlineMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
              overflow: TextOverflow.ellipsis,
            ),
            Text(
              title,
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.outline,
              ),
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ),
      ),
    );
  }
}

class _MetricRow extends StatelessWidget {
  final PublishedContentMetric metric;

  const _MetricRow({required this.metric});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return Card(
      margin: EdgeInsets.only(bottom: defaults.spacing / 2),
      child: Padding(
        padding: defaults.padding,
        child: Row(
          children: [
            Expanded(
              flex: 2,
              child: Text(
                metric.publishedContentId,
                style: theme.textTheme.bodySmall,
                overflow: TextOverflow.ellipsis,
              ),
            ),
            Expanded(
              child: Text(
                '${metric.archivers} archivers',
                style: theme.textTheme.bodyMedium,
                textAlign: TextAlign.end,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
