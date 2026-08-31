import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/community/list.display.item.dart';
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/community.button.publish.dart';
import 'package:retrovibed/community/community.button.delete.dart';
import 'package:retrovibed/community/metrics.dashboard.dart';

void main() {
	group('ListDisplayItem', () {
		testWidgets('owned community shows publish, delete buttons', (tester) async {
			final community = Community(
				id: 'c1',
				accountId: uuidx.min(),
				url: 'https://owned.community.retrovibe.space',
				description: 'my community',
			);

			await tester.pumpApp(ListDisplayItem(community: community));
			await tester.pumpAndSettle();

			expect(find.byType(PublishButton), findsOneWidget);
			expect(find.byType(DeleteButton), findsOneWidget);
			expect(tester.takeException(), isNull);
		});

		testWidgets('non-owned community hides publish, delete buttons', (tester) async {
			final community = Community(
				id: 'c1',
				accountId: 'other-account',
				url: 'https://notmine.community.retrovibe.space',
				description: 'someone else community',
			);

			await tester.pumpApp(ListDisplayItem(community: community));
			await tester.pumpAndSettle();

			expect(find.byType(PublishButton), findsNothing);
			expect(find.byType(DeleteButton), findsNothing);
			expect(tester.takeException(), isNull);
		});

		testWidgets('owned community shows metrics when expanded', (tester) async {
			final community = Community(
				id: 'c1',
				accountId: uuidx.min(),
				url: 'https://owned.community.retrovibe.space',
			);

			await tester.pumpApp(SingleChildScrollView(child: ListDisplayItem(community: community)));
			await tester.pumpAndSettle();
			await tester.tap(find.byType(InkWell).first);
			await tester.pump();

			expect(find.byType(MetricsDashboard), findsOneWidget);
			expect(tester.takeException(), isNull);
		});

		testWidgets('non-owned community hides metrics when expanded', (tester) async {
			final community = Community(
				id: 'c1',
				accountId: 'other-account',
				url: 'https://notmine.community.retrovibe.space',
			);

			await tester.pumpApp(SingleChildScrollView(child: ListDisplayItem(community: community)));
			await tester.pumpAndSettle();
			await tester.tap(find.byType(InkWell).first);
			await tester.pump();

			expect(find.byType(MetricsDashboard), findsNothing);
			expect(tester.takeException(), isNull);
		});

		testWidgets('subscribe button visible for non-owned community', (tester) async {
			final community = Community(
				id: 'c1',
				accountId: 'other-account',
				url: 'https://notmine.community.retrovibe.space',
			);

			await tester.pumpApp(ListDisplayItem(community: community));
			await tester.pumpAndSettle();

			expect(find.byIcon(Icons.add_circle_outline), findsOneWidget);
			expect(tester.takeException(), isNull);
		});

		testWidgets('subscribed non-owned community shows check_circle', (tester) async {
			final community = Community(
				id: 'c1',
				accountId: 'other-account',
				url: 'https://notmine.community.retrovibe.space',
				subscribedAt: '2026-03-20T00:00:00Z',
			);

			await tester.pumpApp(ListDisplayItem(community: community));
			await tester.pumpAndSettle();

			expect(find.byIcon(Icons.check_circle), findsOneWidget);
			expect(find.byIcon(Icons.add_circle_outline), findsNothing);
			expect(tester.takeException(), isNull);
		});

		testWidgets('subscribed owned community shows check_circle', (tester) async {
			final community = Community(
				id: 'c1',
				accountId: uuidx.min(),
				url: 'https://owned.community.retrovibe.space',
				subscribedAt: '2026-03-20T00:00:00Z',
			);

			await tester.pumpApp(ListDisplayItem(community: community));
			await tester.pumpAndSettle();

			expect(find.byIcon(Icons.check_circle), findsOneWidget);
			expect(find.byIcon(Icons.add_circle_outline), findsNothing);
			expect(tester.takeException(), isNull);
		});
	});
}
