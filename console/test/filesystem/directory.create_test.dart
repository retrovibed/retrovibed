import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/filesystem/api.dart' as api;
import 'package:retrovibed/filesystem/directory.create.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<api.FilesystemCreateResponse> _mockCreate(
  api.FilesystemCreateRequest req, {
  List<httpx.Option> options = const [],
}) {
  return Future.value(
    api.FilesystemCreateResponse(media: media.Media(id: 'created-id', description: req.name)),
  );
}

Future<api.FilesystemCreateResponse> _mockCreateFailure(
  api.FilesystemCreateRequest req, {
  List<httpx.Option> options = const [],
}) {
  return Future.error(Exception('create failed'));
}

void main() {
  group('DirectoryCreate', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      await tester.pumpApp(
        DirectoryCreate(
          parent: 'parent-id',
          create: _mockCreate,
          onCreated: (_) {},
          onCancel: () {},
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('submitting a name calls create and onCreated with the resulting media', (
      WidgetTester tester,
    ) async {
      media.Media? created;

      await tester.pumpApp(
        DirectoryCreate(
          parent: 'parent-id',
          create: _mockCreate,
          onCreated: (m) => created = m,
          onCancel: () {},
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), 'New Folder');
      await tester.testTextInput.receiveAction(TextInputAction.done);
      await tester.pumpAndSettle();

      expect(created, isNotNull);
      expect(created?.description, 'New Folder');
      expect(tester.takeException(), isNull);
    });

    testWidgets('submitting an empty name does nothing', (
      WidgetTester tester,
    ) async {
      var created = false;

      await tester.pumpApp(
        DirectoryCreate(
          parent: 'parent-id',
          create: _mockCreate,
          onCreated: (_) => created = true,
          onCancel: () {},
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.check));
      await tester.pumpAndSettle();

      expect(created, isFalse);
      expect(tester.takeException(), isNull);
    });

    testWidgets('a failed create surfaces an error', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        DirectoryCreate(
          parent: 'parent-id',
          create: _mockCreateFailure,
          onCreated: (_) {},
          onCancel: () {},
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), 'New Folder');
      await tester.testTextInput.receiveAction(TextInputAction.done);
      await tester.pumpAndSettle();

      expect(find.text('an unexpected problem has occurred'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping the close icon calls onCancel', (
      WidgetTester tester,
    ) async {
      var cancelled = false;

      await tester.pumpApp(
        DirectoryCreate(
          parent: 'parent-id',
          create: _mockCreate,
          onCreated: (_) {},
          onCancel: () => cancelled = true,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.close));
      await tester.pumpAndSettle();

      expect(cancelled, isTrue);
      expect(tester.takeException(), isNull);
    });
  });
}
