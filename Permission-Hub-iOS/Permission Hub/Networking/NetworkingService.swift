//
//  NetworkingService.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import Foundation

enum PHNetworkingError: Error {
    case invalidUrl
}

extension PHNetworkingError: LocalizedError {

    var errorDescription: String? {
        switch self {
        case .invalidUrl:
            return "unable to construct URL"
        }
    }
}

enum PHNetworkingResult<T> {
    case success(T?)
    case failure(Error)
}

final class NetworkingService {

    // MARK: - Static properties

    static let shared = NetworkingService()

    // MARK: - Private properties

    private let baseUrl = "http://51.15.52.136"

    // MARK: - Endpoints

    // endpoint to poll for notifications
    func getNotifications(completion: @escaping  (_ result: PHNetworkingResult<TransactionNotification>) -> Void) throws {

        // construct URL
        guard let url = URL(string: baseUrl + "/notification-get") else {
            throw PHNetworkingError.invalidUrl
        }

        // make request
        let task = URLSession.shared.dataTask(with: url) { data, response, error in

            guard let response = response as? HTTPURLResponse, response.statusCode == 200 else {
                return
            }

            if let error = error {
                DispatchQueue.main.async {
                    completion(.failure(error))
                }
            }

            if let data = data {

                let decoder = JSONDecoder()
                decoder.dateDecodingStrategy = .iso8601

                do {
//                    let notification = try decoder.decode(TransactionNotification.self, from: data)
                    let data = "{\"transactionID\":\"b80db272b05b9ad007c6833dac68b95ca907594946b2da1929d1f8f95d973b5c\",\"item\":[{\"Item\":\"Access to your Personal Details\",\"Fields\":[\"Name\",\"surname\",\"date of birth\",\"email\",\"BSN\"]},{\"Item\":\"Legal identity (passport)\",\"Fields\":[\"Number\",\"Expiration date\",\"Country of issue\"]},{\"Item\":\"Newsletter\",\"Fields\":[\"Emails for marketing purposes\"]},{\"Item\":\"BankingDetails\",\"Fields\":[\"IBAN\",\"bank\",\"name\"]}],\"title\":\"Provide permission for completing\",\"description\":\"T-mobile monthly plan(unlimited data), 65 euro, iPhone XR 256GB\",\"verification\":[\"digid.nl\",\"planet-blockchain\",\"kvk\"],\"date\":\"2019-04-13T15:51:57+02:00\",\"requesterName\":\"John Smith\",\"RequesterPublicKey\":\"69093eef7426963f2ef0f68fb73e355b7898ddb04a4fad769a96b41ffc824c1c\",\"analysis\":[\"personal data is GDPR protected data\",\"banking details is sensitive data\"]}".data(using: .utf8)!
                    let decoder = JSONDecoder()
                    decoder.dateDecodingStrategy = .iso8601
                    let transaction = try! decoder.decode(TransactionNotification.self, from: data)
                    DispatchQueue.main.async {
                        completion(.success(transaction))
                    }
                } catch {
                    DispatchQueue.main.async {
                        completion(.failure(error))
                    }
                }
            }
        }

        task.resume()
    }

    func respondToTransaction(withId
        id: String,
        isAccepted: Bool,
        completion: @escaping  (_ result: PHNetworkingResult<TransactionNotification>) -> Void) throws {

        // construct URL
        guard let url = URL(string: baseUrl + "/reply-put") else {
            throw PHNetworkingError.invalidUrl
        }

        // create request
        var request = URLRequest(
            url: url,
            cachePolicy: .reloadIgnoringLocalAndRemoteCacheData,
            timeoutInterval: 10)

        let requestBody: [String: Any] = [
            "transactionID": id,
            "accepted": isAccepted,
        ]

        request.httpMethod = "POST"
        request.httpBody = try JSONSerialization.data(withJSONObject: requestBody)

        // make request
        let task = URLSession.shared.dataTask(with: request) { data, response, error in

            guard let response = response as? HTTPURLResponse,
                response.statusCode == 200 else {
                    return
            }

            if let error = error {
                DispatchQueue.main.async {
                    completion(.failure(error))
                }
            }

            DispatchQueue.main.async {
                completion(.success(nil))
            }
        }

        task.resume()
    }
}
