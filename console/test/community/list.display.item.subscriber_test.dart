import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/community/list.display.item.subscriber.dart';
import 'package:retrovibed/community/community.pb.dart';

final _community = Community(
	id: 'c1',
	domain: 'lorem-ipsum-dolor-sit-amet-consectetur-adipiscing-elit-sed-do-eiusmod',
	description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.',
);

void main() {
	group('SubscriberListDisplayItem', () {
		final resolutions = Resolutions.variant();

		testWidgets('renders without overflow', (tester) async {
			final entry = resolutions.currentValue!;
			await tester.pumpApp(
				physicalSize: entry.value,
				ListView(
					children: [
						SubscriberListDisplayItem(community: _community),
					],
				),
			);
			await tester.pumpAndSettle();

			expect(tester.takeException(), isNull);
		}, variant: resolutions);

		testWidgets('renders expanded without overflow', (tester) async {
			final entry = resolutions.currentValue!;
			await tester.pumpApp(
				physicalSize: entry.value,
				ListView(
					children: [
						SubscriberListDisplayItem(community: _community),
					],
				),
			);
			await tester.pumpAndSettle();
			await tester.tap(find.byType(InkWell).first);
			await tester.pump();

			expect(tester.takeException(), isNull);
		}, variant: resolutions);
	});
}
