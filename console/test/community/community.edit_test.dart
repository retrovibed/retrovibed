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
      url: 'https://mysite.community.retrovibe.space',
      description: 'A test community',
    );

    testWidgets('displays url and description fields', (tester) async {
      await tester.pumpApp(CommunityEdit(community: empty, onChange: (_) {}));
      await tester.pumpAndSettle();

      expect(find.text('URL'), findsOneWidget);
      expect(find.text('Description'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows placeholder helper when url is empty', (tester) async {
      await tester.pumpApp(CommunityEdit(community: empty, onChange: (_) {}));
      await tester.pumpAndSettle();

      expect(
        find.text('https://example.community.retrovibe.space'),
        findsOneWidget,
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('helper reflects community url', (tester) async {
      await tester.pumpApp(CommunityEdit(community: populated, onChange: (_) {}));
      await tester.pumpAndSettle();

      expect(
        find.text('https://mysite.community.retrovibe.space'),
        findsAtLeastNWidgets(1),
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('blank url fails validation', (tester) async {
      final formKey = GlobalKey<FormState>();

      await tester.pumpApp(
        Form(
          key: formKey,
          child: CommunityEdit(community: empty, onChange: (_) {}),
        ),
      );
      await tester.pumpAndSettle();

      expect(formKey.currentState!.validate(), isFalse);
      await tester.pumpAndSettle();

      expect(find.text('domain name is required'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('rejects a malformed url', (tester) async {
      final formKey = GlobalKey<FormState>();

      await tester.pumpApp(
        Form(
          key: formKey,
          child: CommunityEdit(community: empty, onChange: (_) {}),
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).first, 'not a url');
      await tester.pump();

      expect(formKey.currentState!.validate(), isFalse);
      await tester.pumpAndSettle();

      expect(find.text('must be a valid url'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('passes validation when a valid url is provided', (tester) async {
      final formKey = GlobalKey<FormState>();

      await tester.pumpApp(
        Form(
          key: formKey,
          child: CommunityEdit(community: populated, onChange: (_) {}),
        ),
      );
      await tester.pumpAndSettle();

      // the field displays just the domain (see 'bare subdomain input fails
      // validation' below), so validation of an unedited field always fails;
      // enter the full absolute url explicitly to exercise the passing case.
      await tester.enterText(
        find.byType(TextFormField).first,
        'https://mysite.community.retrovibe.space',
      );
      await tester.pump();

      expect(formKey.currentState!.validate(), isTrue);
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onChange when url is edited', (tester) async {
      Community? received;

      await tester.pumpApp(
        CommunityEdit(community: empty, onChange: (c) => received = c),
      );
      await tester.pumpAndSettle();

      await tester.enterText(
        find.byType(TextFormField).first,
        'https://newdomain.community.retrovibe.space',
      );
      await tester.pump();

      expect(received, isNotNull);
      expect(received!.url, equals('https://newdomain.community.retrovibe.space'));
    });

    testWidgets('calls onChange with a canonical url when a fully qualified custom url is entered', (tester) async {
      Community? received;

      await tester.pumpApp(
        CommunityEdit(community: empty, onChange: (c) => received = c),
      );
      await tester.pumpAndSettle();

      await tester.enterText(
        find.byType(TextFormField).first,
        'https://custom.example.com',
      );
      await tester.pump();

      expect(received, isNotNull);
      expect(received!.url, equals('https://custom.example.com'));
    });

    testWidgets('calls onChange with a canonical url when a bare subdomain is entered', (tester) async {
      Community? received;

      await tester.pumpApp(
        CommunityEdit(community: empty, onChange: (c) => received = c),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).first, 'mysite');
      await tester.pump();

      expect(received, isNotNull);
      expect(received!.url, equals('https://mysite.community.retrovibe.space'));
    });

    testWidgets('passes validation when a fully qualified custom url is provided', (tester) async {
      final formKey = GlobalKey<FormState>();
      final custom = Community(url: 'https://custom.example.com');

      await tester.pumpApp(
        Form(
          key: formKey,
          child: CommunityEdit(community: custom, onChange: (_) {}),
        ),
      );
      await tester.pumpAndSettle();

      expect(formKey.currentState!.validate(), isTrue);
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('bare subdomain input fails validation even though onChange normalizes it', (tester) async {
      final formKey = GlobalKey<FormState>();

      await tester.pumpApp(
        Form(
          key: formKey,
          child: CommunityEdit(community: empty, onChange: (_) {}),
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).first, 'mysite');
      await tester.pump();

      // NOTE: this documents current behavior. onChange normalizes 'mysite' into
      // a canonical community url via api.communities.canonicaluri, but the
      // validator checks the raw (un-normalized) field value, so it is rejected.
      expect(formKey.currentState!.validate(), isFalse);
      await tester.pumpAndSettle();

      expect(find.text('must be a valid url'), findsOneWidget);
      expect(tester.takeException(), isNull);
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
      final hidden = Community(url: 'https://secret.community.retrovibe.space', hidden: true);

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

      testWidgets('long url', (tester) async {
        final entry = resolutions.currentValue!;

        final longURL = Community(
          url: 'https://a-very-long-domain-name-that-might-cause-overflow.community.retrovibe.space',
          description: 'A description that is also quite long and verbose',
        );

        await tester.pumpApp(
          physicalSize: entry.value,
          SingleChildScrollView(
            child: CommunityEdit(community: longURL, onChange: (_) {}),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: resolutions);
    });
  });
}
