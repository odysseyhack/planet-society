//
//  TransactionWalletViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionWalletViewController: PHTableViewController {

    // MARK: - Life cycle

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
                    self.presentTransactionOverviewViewController(transaction: transaction)
                }
            case .failure(let error):
                print(error)
            }
        }

        // poll endpoint every second
        perform(#selector(pollForNotifications), with: nil, afterDelay: 1)
    }

    // MARK: - Helpers

    private func presentTransactionOverviewViewController(transaction: TransactionNotification) {

        let viewController = TransactionFlowViewController(transaction: transaction)
        navigationController?.pushViewController(viewController, animated: true)
    }
}
