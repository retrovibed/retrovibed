import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta/api.dart' as api;
import 'package:retrovibed/meta/daemon.auto.dart';
import 'package:retrovibed/meta/daemon.mdns.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('EndpointAuto', () {
    testWidgets('attempts the latest daemon on init', (
      WidgetTester tester,
    ) async {
      bool latestCalled = false;
      Future<api.DaemonLookupResponse> mockLatest() async {
        latestCalled = true;
        return api.DaemonLookupResponse(
          daemon: api.Daemon(hostname: 'localhost:9998'),
        );
      }

      Future<api.Daemon> mockConnectable(api.Daemon d) async => d;

      await tester.pumpApp(
        MaterialApp(
          home: EndpointAuto(
            latest: mockLatest,
            connectable: mockConnectable,
            backoff: httpx.Backoff.constant(Duration.zero),
            const SizedBox(),
          ),
        ),
      );

      await tester.pumpAndSettle();

      expect(latestCalled, true);
    });

    testWidgets('shows unknown error on non-retryable exception', (
      WidgetTester tester,
    ) async {
      Future<api.DaemonLookupResponse> mockLatest() async {
        throw Exception('some generic error');
      }

      await tester.pumpApp(
        MaterialApp(
          home: EndpointAuto(
            latest: mockLatest,
            backoff: httpx.Backoff.constant(Duration.zero),
            const SizedBox(),
          ),
        ),
      );

      await tester.pumpAndSettle();

      expect(find.text('an unexpected problem has occurred'), findsOneWidget);
    });

    testWidgets('shows conflict error on http 409', (
      WidgetTester tester,
    ) async {
      Future<api.DaemonLookupResponse> mockLatest() async {
        throw http.Response('', 409);
      }

      await tester.pumpApp(
        MaterialApp(
          home: EndpointAuto(
            latest: mockLatest,
            backoff: httpx.Backoff.constant(Duration.zero),
            const SizedBox(),
          ),
        ),
      );

      await tester.pumpAndSettle();

      expect(
        find.text("you've not yet been approved to access this library"),
        findsOneWidget,
      );
    });

    testWidgets('shows forbidden error on http 403', (
      WidgetTester tester,
    ) async {
      Future<api.DaemonLookupResponse> mockLatest() async {
        throw http.Response('', 403);
      }

      await tester.pumpApp(
        MaterialApp(
          home: EndpointAuto(
            latest: mockLatest,
            backoff: httpx.Backoff.constant(Duration.zero),
            const SizedBox(),
          ),
        ),
      );

      await tester.pumpAndSettle();

      expect(
        find.text(
          "you've attempted to access a service you havent been granted access to yet.",
        ),
        findsOneWidget,
      );
    });

    testWidgets('handles successful daemon connection', (
      WidgetTester tester,
    ) async {
      final daemon = api.Daemon(hostname: 'localhost:8080');
      Future<api.DaemonLookupResponse> mockLatest() async => api.DaemonLookupResponse(daemon: daemon);
      Future<api.Daemon> mockConnectable(api.Daemon d) async => d;

      await tester.pumpApp(
        MaterialApp(
          home: EndpointAuto(
            latest: mockLatest,
            connectable: mockConnectable,
            backoff: httpx.Backoff.constant(Duration.zero),
            const Placeholder(),
          ),
        ),
      );

      await tester.pumpAndSettle();

      expect(find.byType(Placeholder), findsOneWidget);
    });

    testWidgets('retries latest on ClientException and connects', (
      WidgetTester tester,
    ) async {
      int attempts = 0;
      final daemon = api.Daemon(hostname: 'localhost:9998');

      Future<api.DaemonLookupResponse> mockLatest() async {
        attempts++;
        if (attempts == 1) {
          throw http.ClientException(
            'Connection refused',
            Uri.https('localhost:9998', '/meta/d/latest'),
          );
        }
        return api.DaemonLookupResponse(daemon: daemon);
      }

      Future<api.Daemon> mockConnectable(api.Daemon d) async => d;

      await tester.pumpApp(
        ds.LoadingGuard(
          EndpointAuto(
            latest: mockLatest,
            connectable: mockConnectable,
            backoff: httpx.Backoff.constant(Duration.zero),
            const Placeholder(),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(attempts, equals(2));
      expect(find.byType(Placeholder), findsOneWidget);
    });

    testWidgets(
      'shows initial setup when no daemon exists (404 on latest and create)',
      (WidgetTester tester) async {
        Future<api.DaemonLookupResponse> mockLatest() async {
          throw http.Response('', 404);
        }

        Future<api.DaemonCreateResponse> mockCreate(
          api.DaemonCreateRequest r,
        ) async {
          throw http.Response('', 404);
        }

        await tester.pumpApp(
          ds.LoadingGuard(
            EndpointAuto(
              latest: mockLatest,
              create: mockCreate,
              backoff: httpx.Backoff.constant(Duration.zero),
              const SizedBox(),
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(find.byType(InitialSetup), findsOneWidget);
      },
    );

    testWidgets('shows unauthorized error when connectable returns 401', (
      WidgetTester tester,
    ) async {
      final daemon = api.Daemon(hostname: 'localhost:9998');
      Future<api.DaemonLookupResponse> mockLatest() async => api.DaemonLookupResponse(daemon: daemon);
      Future<api.Daemon> mockConnectable(api.Daemon d) async {
        throw http.Response('', 401);
      }

      await tester.pumpApp(
        ds.LoadingGuard(
          EndpointAuto(
            latest: mockLatest,
            connectable: mockConnectable,
            backoff: httpx.Backoff.constant(Duration.zero),
            const SizedBox(),
          ),
        ),
      );

      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(
        find.text(
          "you've attempted to access a system you havent been granted access to yet.",
        ),
        findsOneWidget,
      );
    });

    testWidgets('retries latest on MissingTokenError and connects', (
      WidgetTester tester,
    ) async {
      int attempts = 0;
      final daemon = api.Daemon(hostname: 'localhost:9998');

      Future<api.DaemonLookupResponse> mockLatest() async {
        attempts++;
        if (attempts == 1) throw const httpx.MissingTokenError();
        return api.DaemonLookupResponse(daemon: daemon);
      }

      Future<api.Daemon> mockConnectable(api.Daemon d) async => d;

      await tester.pumpApp(
        ds.LoadingGuard(
          EndpointAuto(
            latest: mockLatest,
            connectable: mockConnectable,
            backoff: httpx.Backoff.constant(Duration.zero),
            const Placeholder(),
          ),
        ),
      );

      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(attempts, equals(2));
      expect(find.byType(Placeholder), findsOneWidget);
    });

    testWidgets('shows unknown error when connectable throws MissingTokenError', (
      WidgetTester tester,
    ) async {
      final daemon = api.Daemon(hostname: 'localhost:9998');
      Future<api.DaemonLookupResponse> mockLatest() async => api.DaemonLookupResponse(daemon: daemon);
      Future<api.Daemon> mockConnectable(api.Daemon d) async {
        throw const httpx.MissingTokenError();
      }

      await tester.pumpApp(
        ds.LoadingGuard(
          EndpointAuto(
            latest: mockLatest,
            connectable: mockConnectable,
            backoff: httpx.Backoff.constant(Duration.zero),
            const SizedBox(),
          ),
        ),
      );

      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(find.text('an unexpected problem has occurred'), findsOneWidget);
    });

    testWidgets(
      'shows forbidden error when connectable returns 403',
      (WidgetTester tester) async {
        final daemon = api.Daemon(hostname: 'localhost:9998');
        Future<api.DaemonLookupResponse> mockLatest() async => api.DaemonLookupResponse(daemon: daemon);
        Future<api.Daemon> mockConnectable(api.Daemon d) async {
          throw http.Response('', 403);
        }

        await tester.pumpApp(
          ds.LoadingGuard(
            EndpointAuto(
              latest: mockLatest,
              connectable: mockConnectable,
              backoff: httpx.Backoff.constant(Duration.zero),
              const SizedBox(),
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(
          find.text(
            "you've attempted to access a service you havent been granted access to yet.",
          ),
          findsOneWidget,
        );
      },
    );

    testWidgets(
      'register flow: connectable 401 then tap Yes connects and shows child',
      (WidgetTester tester) async {
        final daemon = api.Daemon(hostname: 'localhost:9998');
        bool registerCalled = false;
        Future<api.DaemonLookupResponse> mockLatest() async => api.DaemonLookupResponse(daemon: daemon);
        Future<api.Daemon> mockConnectable(api.Daemon d) async {
          throw http.Response('', 401);
        }

        Future<api.Session> mockRegister(
          api.Identity identity, {
          String? host,
        }) async {
          registerCalled = true;
          return api.Session();
        }

        await tester.pumpApp(
          ds.LoadingGuard(
            EndpointAuto(
              latest: mockLatest,
              connectable: mockConnectable,
              register: mockRegister,
              backoff: httpx.Backoff.constant(Duration.zero),
              const Placeholder(),
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(find.text('Yes'), findsOneWidget);
        await tester.tap(find.text('Yes'));
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
        expect(registerCalled, true);
        expect(find.byType(Placeholder), findsOneWidget);
      },
    );

    testWidgets(
      'register flow: connectable 401 then tap Yes with register 409 shows conflict error',
      (WidgetTester tester) async {
        final daemon = api.Daemon(hostname: 'localhost:9998');
        Future<api.DaemonLookupResponse> mockLatest() async => api.DaemonLookupResponse(daemon: daemon);
        Future<api.Daemon> mockConnectable(api.Daemon d) async {
          throw http.Response('', 401);
        }

        Future<api.Session> mockRegister(
          api.Identity identity, {
          String? host,
        }) async {
          throw http.Response('', 409);
        }

        await tester.pumpApp(
          ds.LoadingGuard(
            EndpointAuto(
              latest: mockLatest,
              connectable: mockConnectable,
              register: mockRegister,
              backoff: httpx.Backoff.constant(Duration.zero),
              const Placeholder(),
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsOneWidget);
        await tester.tap(find.text('Yes'));
        await tester.pumpAndSettle();

        expect(
          find.text("you've not yet been approved to access this library"),
          findsOneWidget,
        );
      },
    );

    testWidgets('shows NoLocalService on ENOROUTE socket error', (
      WidgetTester tester,
    ) async {
      Future<api.DaemonLookupResponse> mockLatest() async {
        throw SocketException('', osError: OSError('', 113));
      }

      await tester.pumpApp(
        ds.LoadingGuard(
          EndpointAuto(
            latest: mockLatest,
            backoff: httpx.Backoff.constant(Duration.zero),
            const SizedBox(),
          ),
        ),
      );

      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(find.byType(NoLocalService), findsOneWidget);
    });

    testWidgets('creates daemon when latest returns 404 and connects', (
      WidgetTester tester,
    ) async {
      final daemon = api.Daemon(hostname: 'localhost:9998');
      Future<api.DaemonLookupResponse> mockLatest() async {
        throw http.Response('', 404);
      }

      Future<api.DaemonCreateResponse> mockCreate(
        api.DaemonCreateRequest r,
      ) async {
        return api.DaemonCreateResponse(daemon: daemon);
      }

      Future<api.Daemon> mockConnectable(api.Daemon d) async => d;

      await tester.pumpApp(
        MaterialApp(
          home: ds.LoadingGuard(
            EndpointAuto(
              latest: mockLatest,
              create: mockCreate,
              connectable: mockConnectable,
              backoff: httpx.Backoff.constant(Duration.zero),
              const Placeholder(),
            ),
          ),
        ),
      );

      await tester.pumpAndSettle();

      expect(find.byType(Placeholder), findsOneWidget);
    });

    for (final entry in Resolutions.all.entries) {
      testWidgets(
        'shows NoLocalService on offline socket error (loose) at ${entry.key}',
        (WidgetTester tester) async {
          Future<api.DaemonLookupResponse> mockLatest() async {
            throw SocketException('', osError: OSError('', 111));
          }

          await tester.pumpApp(
            physicalSize: entry.value,
            ds.LoadingGuard(
              EndpointAuto(
                latest: mockLatest,
                backoff: httpx.Backoff.constant(Duration.zero),
                const SizedBox(),
              ),
            ),
            fit: FlexFit.loose,
          );
          await tester.pumpAndSettle();
          expect(find.byType(NoLocalService), findsOneWidget);
        },
      );

      testWidgets(
        'shows NoLocalService on offline socket error (tight) at ${entry.key}',
        (WidgetTester tester) async {
          Future<api.DaemonLookupResponse> mockLatest() async {
            throw SocketException('', osError: OSError('', 111));
          }

          await tester.pumpApp(
            physicalSize: entry.value,
            ds.LoadingGuard(
              EndpointAuto(
                latest: mockLatest,
                backoff: httpx.Backoff.constant(Duration.zero),
                const SizedBox(),
              ),
            ),
            fit: FlexFit.tight,
          );
          await tester.pumpAndSettle();
          expect(find.byType(NoLocalService), findsOneWidget);
        },
      );
    }

    for (final entry in Resolutions.all.entries) {
      testWidgets(
        'shows NoLocalService on dns resolution failure (loose) at ${entry.key}',
        (WidgetTester tester) async {
          Future<api.DaemonLookupResponse> mockLatest() async {
            throw SocketException('', osError: OSError('', -2));
          }

          await tester.pumpApp(
            physicalSize: entry.value,
            ds.LoadingGuard(
              EndpointAuto(
                latest: mockLatest,
                backoff: httpx.Backoff.constant(Duration.zero),
                const SizedBox(),
              ),
            ),
            fit: FlexFit.loose,
          );
          await tester.pumpAndSettle();
          expect(find.byType(NoLocalService), findsOneWidget);
        },
      );

      testWidgets(
        'shows NoLocalService on dns resolution failure (tight) at ${entry.key}',
        (WidgetTester tester) async {
          Future<api.DaemonLookupResponse> mockLatest() async {
            throw SocketException('', osError: OSError('', -2));
          }

          await tester.pumpApp(
            physicalSize: entry.value,
            ds.LoadingGuard(
              EndpointAuto(
                latest: mockLatest,
                backoff: httpx.Backoff.constant(Duration.zero),
                const SizedBox(),
              ),
            ),
            fit: FlexFit.tight,
          );
          await tester.pumpAndSettle();
          expect(find.byType(NoLocalService), findsOneWidget);
        },
      );

      testWidgets(
        'shows NoLocalService on dns resolution failure (MaterialApp) at ${entry.key}',
        (WidgetTester tester) async {
          Future<api.DaemonLookupResponse> mockLatest() async {
            throw SocketException('', osError: OSError('', -2));
          }

          await tester.pumpApp(
            physicalSize: entry.value,
            Scaffold(
              body: ds.LoadingGuard(
                EndpointAuto(
                  latest: mockLatest,
                  backoff: httpx.Backoff.constant(Duration.zero),
                  const SizedBox(),
                ),
              ),
            ),
          );
          await tester.pumpAndSettle();
          expect(find.byType(NoLocalService), findsOneWidget);
        },
      );
    }

    testWidgets('does not show MDNSDiscovery while loading', (
      WidgetTester tester,
    ) async {
      Future<api.DaemonLookupResponse> mockLatest() async {
        await Future.delayed(const Duration(seconds: 1));
        return api.DaemonLookupResponse(
          daemon: api.Daemon(hostname: 'localhost:9998'),
        );
      }

      Future<api.Daemon> mockConnectable(api.Daemon d) async => d;

      await tester.pumpApp(
        MaterialApp(
          home: ds.LoadingGuard(
            EndpointAuto(
              latest: mockLatest,
              connectable: mockConnectable,
              backoff: httpx.Backoff.constant(Duration.zero),
              const SizedBox(),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // still loading — MDNSDiscovery must not appear until loading completes
      expect(find.byType(MDNSDiscovery), findsNothing);
      await tester.pumpAndSettle();
    });

    testWidgets('shows loading indicator during daemon lookup', (
      WidgetTester tester,
    ) async {
      Future<api.DaemonLookupResponse> mockLatest() async {
        await Future.delayed(const Duration(seconds: 1));
        return api.DaemonLookupResponse(
          daemon: api.Daemon(hostname: 'localhost:9998'),
        );
      }

      Future<api.Daemon> mockConnectable(api.Daemon d) async => d;

      await tester.pumpApp(
        MaterialApp(
          home: ds.LoadingGuard(
            EndpointAuto(
              latest: mockLatest,
              connectable: mockConnectable,
              backoff: httpx.Backoff.constant(Duration.zero),
              const SizedBox(),
            ),
          ),
        ),
      );
      await tester.pump();
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsOneWidget);

      await tester.pumpAndSettle();

      expect(find.byType(CircularProgressIndicator), findsNothing);
    });
  });
}
