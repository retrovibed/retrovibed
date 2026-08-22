import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/remote/connect.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<Stream<meta.Daemon>> _noopDaemonDiscover({List<httpx.Option> options = const []}) async {
  return const Stream<meta.Daemon>.empty();
}

void main() {
  testWidgets('unmounting Connect does not throw during dispose', (tester) async {
    bool visible = true;
    late StateSetter setLocalState;

    await tester.pumpApp(
      StatefulBuilder(
        builder: (context, setState) {
          setLocalState = setState;
          return visible
              ? Connect(
                  search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
                )
              : const SizedBox();
        },
      ),
    );
    await tester.pump();

    setLocalState(() => visible = false);
    await tester.pump();

    expect(tester.takeException(), isNull);
  });

  testWidgets('selecting the local device is refused without calling connect', (tester) async {
    bool connectCalled = false;

    await tester.pumpApp(
      Connect(
        search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
        daemonDiscover: _noopDaemonDiscover,
        connect: ({required String host, List<httpx.Option> options = const []}) async {
          connectCalled = true;
          throw StateError('connect should not be called for the local device');
        },
      ),
    );
    await tester.pumpN(5);

    authn.AuthedEndpoint.daemon(tester.element(find.byType(meta.DaemonDropdown))).value = meta.Daemon(
      hostname: "localhost:9998",
    );
    await tester.pumpN(5);

    expect(find.text("you do not have permission to remotely control this device"), findsOneWidget);
    expect(connectCalled, isFalse);
    expect(tester.takeException(), isNull);
  });

  testWidgets('a forbidden response from connect disables remote control', (tester) async {
    await tester.pumpApp(
      Connect(
        search: ValueNotifier(media.MediaSearchState(next: media.MediaSearchRequest())),
        daemonDiscover: _noopDaemonDiscover,
        connect: ({required String host, List<httpx.Option> options = const []}) async {
          throw http.Response('', 403);
        },
      ),
    );
    await tester.pumpN(5);

    authn.AuthedEndpoint.daemon(tester.element(find.byType(meta.DaemonDropdown))).value = meta.Daemon(
      hostname: "example.remote:1234",
    );
    await tester.pumpN(5);

    expect(find.text("remote control is disabled on this device"), findsOneWidget);
    expect(tester.takeException(), isNull);
  });
}
