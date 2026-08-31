import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/community.dart';
import 'package:retrovibed/community/metrics.dashboard.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _now = DateTime(2024, 6, 15);
final _segments = [
  timex.Range(DateTime(_now.year, _now.month - 1, _now.day), _now),
  timex.Range(DateTime(_now.year, _now.month - 3, _now.day), _now),
  timex.Range(DateTime(_now.year - 1, _now.month, _now.day), _now),
  timex.Range(DateTime(_now.year - 3, _now.month, _now.day), _now),
];

Future<MetricsSyncProgress> mockSync(
  String id, {
  List<httpx.Option> options = const [],
}) async {
  return MetricsSyncProgress();
}

Future<CommunityMetricsResponse> mockMetrics(
  String id, {
  required DateTime startDate,
  required DateTime endDate,
  List<httpx.Option> options = const [],
}) async {
  return CommunityMetricsResponse(
    summary: CommunityMetric(subscribers: 42),
    totalArchivers: 10,
    items: [
      PublishedContentMetric(publishedContentId: 'content-1', archivers: 5),
      PublishedContentMetric(publishedContentId: 'content-2', archivers: 3),
    ],
  );
}

Future<CommunityMetricsResponse> mockEmptyMetrics(
  String id, {
  required DateTime startDate,
  required DateTime endDate,
  List<httpx.Option> options = const [],
}) async {
  return CommunityMetricsResponse();
}

Community testCommunity() {
  return Community(
    id: 'test-community-id',
    url: 'https://testdomain.community.retrovibe.space',
    description: 'A test community',
    createdAt: '2024-01-15T14:30:00Z',
  );
}

final _resolutions = Resolutions.variant();

void main() {
  group('MetricsDashboard', () {
    group('resolutions - empty state', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          MetricsDashboard(
            community: testCommunity(),
            segments: _segments,
            apicommunitysync: mockSync,
            apicommunitymetrics: mockEmptyMetrics,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Metrics'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('resolutions - with metrics data', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          MetricsDashboard(
            community: testCommunity(),
            segments: _segments,
            apicommunitysync: mockSync,
            apicommunitymetrics: mockMetrics,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Metrics'), findsOneWidget);
        expect(find.text('42'), findsOneWidget);
        expect(find.text('Subscribers'), findsOneWidget);
        expect(find.text('10'), findsOneWidget);
        expect(find.text('Archivers'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });
  });
}
