import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/community/community.detail.dart';
import 'package:retrovibed/community/community.pb.dart';

void main() {
  group('CommunityDetail', () {
    testWidgets('displays domain and URL', (WidgetTester tester) async {
      final community = Community(
        domain: 'testdomain',
        description: '',
        createdAt: '2024-01-15T14:30:00Z',
      );

      await tester.pumpApp(CommunityDetail(community: community));
      await tester.pumpAndSettle();

      expect(find.text('testdomain'), findsOneWidget);
      expect(
        find.text('https://testdomain.community.retrovibe.space'),
        findsOneWidget,
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays description when not empty', (
      WidgetTester tester,
    ) async {
      final community = Community(
        domain: 'example',
        description: 'A test community',
        createdAt: '2024-01-15T14:30:00Z',
      );

      await tester.pumpApp(CommunityDetail(community: community));
      await tester.pumpAndSettle();

      expect(find.text('A test community'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('hides description when empty', (WidgetTester tester) async {
      final community = Community(
        domain: 'example',
        description: '',
        createdAt: '2024-01-15T14:30:00Z',
      );

      await tester.pumpApp(CommunityDetail(community: community));
      await tester.pumpAndSettle();

      // Only domain, URL, and timestamp should be present — no description text
      final column = tester.widget<Column>(find.byType(Column).first);
      // 3 children: domain Text, URL SelectableText, Timestamp
      expect(column.children.length, equals(3));
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays created timestamp', (WidgetTester tester) async {
      final community = Community(
        domain: 'example',
        description: '',
        createdAt: '2024-01-15T14:30:00Z',
      );

      await tester.pumpApp(CommunityDetail(community: community));
      await tester.pumpAndSettle();

      expect(find.text('Created: '), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows lock icon when hidden', (WidgetTester tester) async {
      final community = Community(
        domain: 'secret',
        description: '',
        createdAt: '2024-01-15T14:30:00Z',
        hidden: true,
      );

      await tester.pumpApp(CommunityDetail(community: community));
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.lock), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('hides lock icon when not hidden', (WidgetTester tester) async {
      final community = Community(
        domain: 'public',
        description: '',
        createdAt: '2024-01-15T14:30:00Z',
      );

      await tester.pumpApp(CommunityDetail(community: community));
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.lock), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without error in Expanded within Row', (
      WidgetTester tester,
    ) async {
      final community = Community(
        domain: 'a-very-long-domain-name-that-might-overflow',
        description: 'A description that is also quite long and verbose',
        createdAt: '2024-01-15T14:30:00Z',
      );

      await tester.pumpApp(
        Row(
          children: [
            Expanded(child: CommunityDetail(community: community)),
            SizedBox(width: 50),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });
}
