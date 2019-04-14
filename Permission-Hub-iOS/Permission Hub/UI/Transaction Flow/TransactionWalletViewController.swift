//
//  TransactionWalletViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionWalletViewController: PHTableViewController {

    // MARK: - Initialization

    init() {

        let items: [PHTableViewViewCellType] = [
//            PHTableViewViewCellType.plugin(image: <#T##UIImage?#>, text: <#T##String#>)
        ]

        super.init(items: items)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        navigationItem.title = "My Permission Wallet"
    }

    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(animated)

        // populate view with mock data
        let data = "{\"transactionID\":\"b80db272b05b9ad007c6833dac68b95ca907594946b2da1929d1f8f95d973b5c\",\"item\":[{\"Item\":\"Access to your Personal Details\",\"Fields\":[\"Name\",\"Surname\",\"Date of birth\",\"Email\",\"BSN\"]},{\"Item\":\"Legal identity (passport)\",\"Fields\":[\"Number\",\"Expiration date\",\"Country of issue\"]},{\"Item\":\"Newsletter\",\"Fields\":[\"Email address for marketing purposes\"]},{\"Item\":\"Payment information\",\"Fields\":[\"IBAN number\",\"Bank name\",\"Payment details\"]},{\"Item\":\"Subscription contract 24 months\",\"Fields\":[\"Read and accept the terms\"]}],\"title\":\"Provide permission for completing\",\"description\":\"T-mobile monthly plan(unlimited data), 65 euro, iPhone XR 256GB\",\"verification\":[\"digid.nl\",\"planet-blockchain\",\"kvk\"],\"date\":\"2019-04-13T15:51:57+02:00\",\"requesterName\":\"John Smith\",\"RequesterPublicKey\":\"69093eef7426963f2ef0f68fb73e355b7898ddb04a4fad769a96b41ffc824c1c\",\"analysis\":[\"personal data is GDPR protected data\",\"banking details is sensitive data\"]}".data(using: .utf8)!
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        let transaction = try! decoder.decode(TransactionNotification.self, from: data)
        presentTransactionFlowViewController(transaction: transaction)
    }

    override func viewWillAppear(_ animated: Bool) {
        super.viewWillAppear(animated)

        pollForNotifications()
    }

    // MARK: - Networking

    @objc private func pollForNotifications() {

        let service = NetworkingService.shared
        try! service.getNotifications { result in

            switch result {
            case .success(let transactionOrNil):
                if let transaction = transactionOrNil {
                    self.presentTransactionFlowViewController(transaction: transaction)
                }
            case .failure(let error):
                print(error)
            }
        }

        // poll endpoint every second
        perform(#selector(pollForNotifications), with: nil, afterDelay: 1)
    }

    // MARK: - Helpers

    private func presentTransactionFlowViewController(transaction: TransactionNotification) {

        let viewController = TransactionFlowViewController(transaction: transaction)
        let navigationController = UINavigationController(rootViewController: viewController)
        present(navigationController, animated: true)
    }
}
