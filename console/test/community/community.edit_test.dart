import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/community/community.edit.dart';
import 'package:retrovibed/community/community.pb.dart';

void main() {
  group('CommunityEdit', () {
    final empty = Community();
    final populated = Community(
      id: 'abc',
      accountId: 'acc1',
      domain: 'mysite',
      description: 'A test community',
    );

    testWidgets('displays domain and description fields', (tester) async {
      await tester.pumpApp(CommunityEdit(community: empty, onChange: (_) {}));
      await tester.pumpAndSettle();

      expect(find.text('Domain'), findsOneWidget);
      expect(find.text('Description'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows default helper URL when domain is empty', (tester) async {
      await tester.pumpApp(CommunityEdit(community: empty, onChange: (_) {}));
      await tester.pumpAndSettle();

      expect(
        find.text('https://example.community.retrovibe.space'),
        findsOneWidget,
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('helper URL reflects community domain', (tester) async {
      await tester.pumpApp(CommunityEdit(community: populated, onChange: (_) {}));
      await tester.pumpAndSettle();

      expect(
        find.text('https://mysite.community.retrovibe.space'),
        findsOneWidget,
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('validates that domain is required', (tester) async {
      final formKey = GlobalKey<FormState>();

      await tester.pumpApp(
        Form(
          key: formKey,
          child: CommunityEdit(community: empty, onChange: (_) {}),
        ),
      );
      await tester.pumpAndSettle();

      formKey.currentState!.validate();
      await tester.pumpAndSettle();

      expect(find.text('Domain is required'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('passes validation when domain is provided', (tester) async {
      final formKey = GlobalKey<FormState>();

      await tester.pumpApp(
        Form(
          key: formKey,
          child: CommunityEdit(community: populated, onChange: (_) {}),
        ),
      );
      await tester.pumpAndSettle();

      expect(formKey.currentState!.validate(), isTrue);
      await tester.pumpAndSettle();

      expect(find.text('Domain is required'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onChange when domain is edited', (tester) async {
      Community? received;

      await tester.pumpApp(
        CommunityEdit(community: empty, onChange: (c) => received = c),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).first, 'newdomain');
      await tester.pump();

      expect(received, isNotNull);
      expect(received!.domain, equals('newdomain'));
    });

    testWidgets('calls onChange when description is edited', (tester) async {
      Community? received;

      await tester.pumpApp(
        CommunityEdit(community: populated, onChange: (c) => received = c),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).at(1), 'updated desc');
      await tester.pump();

      expect(received, isNotNull);
      expect(received!.description, equals('updated desc'));
    });

    testWidgets('displays visibility selector defaulting to public', (tester) async {
      await tester.pumpApp(CommunityEdit(community: empty, onChange: (_) {}));
      await tester.pumpAndSettle();

      expect(find.text('Visibility'), findsOneWidget);
      expect(find.text('Public'), findsOneWidget);
      expect(find.text('Private'), findsOneWidget);
      expect(find.text('Discoverable via search.'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('visibility selector reflects hidden community', (tester) async {
      final hidden = Community(domain: 'secret', hidden: true);

      await tester.pumpApp(CommunityEdit(community: hidden, onChange: (_) {}));
      await tester.pumpAndSettle();

      expect(find.text('Hidden from search for other accounts.'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onChange when visibility toggled to private', (tester) async {
      Community? received;

      await tester.pumpApp(
        CommunityEdit(community: empty, onChange: (c) => received = c),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Private'));
      await tester.pumpAndSettle();

      expect(received, isNotNull);
      expect(received!.hidden, isTrue);
    });

    group('renders without overflow', () {
      final resolutions = Resolutions.variant();

      testWidgets('populated community', (tester) async {
        final entry = resolutions.currentValue!;

        await tester.pumpApp(
          physicalSize: entry.value,
          SingleChildScrollView(
            child: CommunityEdit(community: populated, onChange: (_) {}),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: resolutions);

      testWidgets('long domain', (tester) async {
        final entry = resolutions.currentValue!;

        final longDomain = Community(
          domain: 'a-very-long-domain-name-that-might-cause-overflow',
          description: 'A description that is also quite long and verbose',
        );

        await tester.pumpApp(
          physicalSize: entry.value,
          SingleChildScrollView(
            child: CommunityEdit(community: longDomain, onChange: (_) {}),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: resolutions);
    });
  });
}
